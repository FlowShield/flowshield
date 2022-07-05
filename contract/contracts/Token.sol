// Solidity files have to start with this pragma.
// It will be used by the Solidity compiler to validate its version.
pragma solidity ^0.8.15;

// This is the main building block for smart contracts.
contract Token {
    // Some string type variables to identify the token.
    string public name = "My Hardhat Token";
    string public symbol = "MBT";
    uint8 public decimals = 2;


    uint256 public totalSupply = 1000000;

    // A mapping is a key/value map. Here we store each account balance.
    mapping(address => uint256) private balances;
    mapping(address => mapping(address=>uint256)) private allowances;

    event Transfer(
        address indexed from,
        address indexed to,
        uint256 value
    );
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );
    constructor() {
        // The totalSupply is assigned to transaction sender, which is the account
        // that is deploying the contract.
        balances[msg.sender] = totalSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return balances[account];
    }

    /**
    * @dev Transfer token for a specified address
    * @param to The address to transfer to.
    * @param value The amount to be transferred.
    */
    function transfer(address to, uint256 value) public returns (bool) {
        require(balances[msg.sender] >= value);

        balances[msg.sender] -= value;
        balances[to] += value;
        emit Transfer(msg.sender, to, value);
        return true;
    }
    /**
    * @dev Transfer tokens from one address to another
    * @param from address The address which you want to send tokens from
    * @param to address The address which you want to transfer to
    * @param value uint256 the amount of tokens to be transferred
    */
    function transferFrom(address from,address to,uint256 value)public returns (bool){
        require(balances[from] >= value );
        require(allowances[from][msg.sender] >= value);

        balances[from] -= value;
        allowances[from][msg.sender] -= value;
        balances[to] += value;
        emit Transfer(from, to, value);
        return true;
    }
    function approve(address spender, uint256 value)public returns (bool){
        allowances[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    function allowance(address owner, address spender)public view returns (uint256){
        return allowances[owner][spender];
    }
   

    //Stake amount
    //uint256 public stakeAmount = 1000;

    

    // // A mapping is a key/value map. Here we store each staked user.
    // mapping(address => uint256) stakeList;

    // /**
    //  * 是否已经质押
    //  */
    // function isStake() external view returns (bool) {
    //     return stakeList[msg.sender] != 0;
    // }
    
    // /**
    //  * 质押
    //  */
    // function stake() external {
    //     require(balances[msg.sender] >= stakeAmount, "Not enough tokens");
        
    //     balances[msg.sender] -= stakeAmount;
    //     balances[owner] += stakeAmount;
    //     stakeList[msg.sender] = stakeAmount;
    // }
    
}