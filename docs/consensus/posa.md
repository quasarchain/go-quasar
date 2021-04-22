# posa

This document gives a detailed description of the posa consensus algorithm and system contract process.

## Introduction to Consensus Algorithm

The posa algorithm is modified based on the clique algorithm, and the system contract is added to realize the related functions of Validator pledge/access, which can incentivize the Validator (obtain handling fees and hsct token rewards) and punish (confiscate the current income, remove the validator list) ).

System function (golang consensus code function):

1. In the first block, pass parameters (Admin, Premint) to initialize the system contract

2. At the end of the block period (`number%Epoch==0`):
    * Get the current top TopValidators through the system contract getTopValidators, and fill them in the extraData field
    * Call the system contract updateActiveValidatorSet to update the list of Validators currently active in the contract
    * Call the system contract reduceMissedBlocksCounter to try to reduce the number of validator errors and avoid being accidentally removed from the validator list due to the accumulation of errors

3. Only replace the validator list with a new validator list at the first block of the block cycle. At this time, the new validator can produce blocks, and the change of the validator list in the contract will only take effect in the next block cycle.

4. When an out of turn block appears, and the validator that should have produced the block has not produced a block recently, the validator calls the punish interface of the system contract to punish the validator. If the number of validator errors reaches punishThreshold (default 10), the current profit will be confiscated. When the removeThreshold (default 30) is reached, the validator list is kicked out and the status is set to Jailed.


Contract function:

1. Validator contract:
    * User pledges/adds HB to become validator
    * The user redeems HB and exits the validator list
    * Redemption block proceeds

2. Proposal contract:
    * The user creates a proposal and applies to become a validator
    * The validator votes on the proposal

3. Punish contract: mainly called by the system to punish the validator of miss block

4. HSCT Token contract: hsct erc20 contract, pre-mining 25 million to the set pre-mining address. The maximum number of tokens is 100 million. Will automatically calculate the corresponding token rewards to the validator according to the ratio

## User becomes validator process

1. Call the Proposal contract `createProposal` interface to create a proposal, the proposal ID can be obtained in the event

2. Validator votes on the contract

3. After the vote is passed, the user calls the `stake` interface of the Validators contract and pledges at least 32 HBs to become a Validator candidate. If the pledge amount ranks in the top 21, add it to the top validator list.

4. During the block epoch, the system calls the Validators contract to obtain the top validator list, writes it into the extra data, and updates the validator list in the next block. At this time, the new validator can generate blocks.

## Parameter configuration

Consensus parameters are mainly divided into two parts:

1. The block time, block period, etc. can be configured in the genesis block.

2. To configure the contract parameters, you need to modify the contract source code, and set the source code of the corresponding system contract in the genesis block.

### Genesis block configuration

The following must be configured in the genesis block:

-Period: block time interval, set to 0 means that blocks will be generated only when there is a transaction

-Epoch: Update the block number interval of validators, the system updates the validators list in the first block of the block cycle

-Admin: The address of validators contract administrator, you can update the exchange rate of block reward hsct

-Premint: hsct pre-mining address, 25 million hsct tokens will be automatically pre-mined to this address

-Account settings, please set the corresponding code of the system contract in the genesis block (**deployedCode**):
    -Validators(0x000000000000000000000000000000000000f000): Validator contract
    -Punish(0x000000000000000000000000000000000000f001): Punishment contract
    -Proposal(0x000000000000000000000000000000000000f002): Proposal contract
    -HSCTToken(0x000000000000000000000000000000000000f003): hsct token contract



### System contract parameters

General parameters:

-MaxValidators(21): The maximum number of validators activated by default

-StakingLockPeriod(100): The block interval from the validator application to redeem the pledged hb to the actual redemption of the pledged hb.

-RestakingLockPeriodUnstaked(200): When the validator redeems the pledged hb and exits the validator list, it wants to join the block interval that needs to be waited for the next time (the exit block is the block number during the redemption operation).

-RestakingLockPeriodJailed(300):
    -When the validator is removed from the validator list due to the system's penalty such as disconnection, it wants to re-add the block interval that needs to be waited (the exit block is the block number during the removal operation)
    -The block interval should be greater than the setting value of RestakingLockPeriodUnstaked

-WithdrawProfitPeriod(100): The smallest interval block size between validators for successive redemption proceeds. When the system slashes the validator, the current unredeemed proceeds will be cleared to zero and distributed equally to other validators.

-MinialStakingCoin(32 ether): The minimum amount of staking hb to be a validator candidate



Punish contract:

-punishThreshold(10): confiscation of the current income error threshold; when the validator is offline and the number of unproduced blocks reaches this threshold, the validator’s current income will be confiscated and distributed equally to other validators

-removeThreshold(30): Remove the error threshold of the validator list; when the validator is offline and the number of unproduced blocks reaches this threshold, its current income will be confiscated and removed from the validator list

-decreaseRate(4): the ratio of validator error clearing. After each epoch, the system will automatically send system transactions to reduce the number of validator errors according to a certain percentage (to prevent the error ratio from being small, but the value is always accumulated to cause punishment). Reduction rules: If the validator error times do not exceed removeThreshold / decreaseRate, no action will be taken; otherwise, the number of errors will be reduced by removeThreshold / decreaseRate


Proposal contract:

-proposalLastingPeriod(7 days): The existence time of the proposal. After this time period, the proposal is invalid if it is not passed.

## System contract user interface

This section introduces the interfaces that can be called externally.

### Proposal

If a user wants to become a validator, he or another person must create a proposal and apply to become a validator. The currently activated validator can vote on the proposal. When the number of people agreeing exceeds one and a half, the proposal is approved (permanently valid), and the user can later become a validator candidate by staking hb.

#### createProposal

Any user creates a proposal (the proposal lasts for 7 days by default, and the proposal needs to be recreated if the proposal is not passed within 7 days).

```solidity
# dst: The address of a candidate to become a validator
# details: detailed description of validator candidates (optional, length should not be greater than 3000)
createProposal(address dst, string calldata details)



# Transaction log
# id: Proposal id, which can be used for voting
# proposer: proposer
# dst: validator candidate address
# time: Proposal time
event LogCreateProposal(
    bytes32 indexed id,
    address indexed proposer,
    address indexed dst,
    uint256 time
);
```

#### voteProposal

The current validator votes on the proposal. When half of the votes were agreed, the proposal passed

```solidity
# id: proposal id
# auth: Do you agree to the proposal
voteProposal(bytes32 id, bool auth)


# Transaction generation log
# Proposal approval log
# id: proposal id
# dst: validator candidate address
# time: Pass time
event LogPassProposal(
    bytes32 indexed id,
    address indexed dst,
    uint256 time
);
# Proposal failed log (more than half of them disagree)
# id: proposal id
# dst: validator candidate address
# time: Proposal failed time
event LogRejectProposal(
    bytes32 indexed id,
    address indexed dst,
    uint256 time
);
```
### Validators

The validator/admin calls the contract to perform related operations such as pledge, redemption of deposits, and redemption of proceeds.

The administrator calls the interface:

```solidity
# The administrator updates the hsct block reward exchange ratio
# multi_: multiplier
# divisor_: divisor
# Actual proportional hsct quantity = block handling fee*multi_/divisor_
changeDec(uint256 multi_, uint256 divisor_)


# Transaction log
# multi_: multiplier
# divisor_: divisor
event LogChangeDec(uint256 newMulti, uint256 newDivisor);

```

validator call interface:

```solidity
# validator pledge/add hb
# feeAddr: benefit address
# moniker: name, the length is not greater than 70
# identity: Identity information, the length is not greater than 3000
# website: Website information, the length is not greater than 140
# email: Email message, the length is not more than 140
# details: Detailed information, the length is not greater than 280
stake(
    address payable feeAddr,
    string calldata moniker,
    string calldata identity,
    string calldata website,
    string calldata email,
    string calldata details
)

# Transaction log
# Become a validator for the first time
# val: validator address
# fee: beneficiary address
# staking: The number of hb pledged
# time: Transaction time
event LogCreateValidator(
    address indexed val,
    address indexed fee,
    uint256 staking,
    uint256 time
);
# Additional staking
# val: validator address
# addAmount: Additional pledge amount
event LogAddStake(address indexed val, uint256 addAmount);
# Become a top validator
# val: validator address
# time: Transaction time
event LogAddToTopValidators(address indexed val, uint256 time);



# Edit validator information
# feeAddr: benefit address
# moniker: name, the length is not greater than 70
# identity: Identity information, the length is not greater than 3000
# website: Website information, the length is not greater than 140
# email: Email message, the length is not more than 140
# details: Detailed information, the length is not greater than 280
editValidator(
    address payable feeAddr,
    string calldata moniker,
    string calldata identity,
    string calldata website,
    string calldata email,
    string calldata details
)
# Transaction log
# val: validator address
# fee: updated benefit address
# time: Update time
event LogEditValidator(
    address indexed val,
    address indexed fee,
    uint256 time
);


# Re-pledge, this method can only be called after the validator is jailed by the system and the latter redeems the deposit and exits the validator list
restake()
# Transaction log
# val: validator address
# staking: Staking amount
# time: Transaction time
event LogRestake(address indexed val, uint256 staking, uint256 time);
# Become a top validator
# val: validator address
# time: Transaction time
event LogAddToTopValidators(address indexed val, uint256 time);



# validatorApply to exit the validator list
# Note: The deposit will not be sent to the validator immediately. The deposit can be redeemed by calling withdrawStaking() after the time set by the system (100 blocks)
unstake()
# Transaction log
# val: validator address
# time: Transaction time
event LogUnstake(address indexed val, uint256 time);

# Redeem deposit
withdrawStaking()
# Transaction log
# val: validator address
# amount: Deposit amount
# time: Transaction time
event LogWithdrawStaking(address indexed val, uint256 amount, uint256 time);



# Redeem the block revenue (including hb and hsct), this method is called by the revenue address, and which validator's revenue needs to be redeemed when calling
# validator: validator address
withdrawProfits(address validator)
# Transaction event
# val: address of validator
# fee: beneficiary's address
# hb: hb income amount
# hsct: hsct income amount
event LogWithdrawProfits(
    address indexed val,
    address indexed fee,
    uint256 hb,
    uint256 hsct
);



# View the list of currently activated validators, the list of validators that can currently produce blocks
getActiveValidators() returns (address[] memory)
# View the current top validator list, the current validator list with the highest pledge amount, and it will be activated in the next cycle
getTopValidators() returns (address[] memory)
```