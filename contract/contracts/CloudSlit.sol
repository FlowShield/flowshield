// Solidity files have to start with this pragma.
// It will be used by the Solidity compiler to validate its version.
pragma solidity ^0.8.15;

// This is the main building block for smart contracts.
contract CloudSlit {
    // Some string type variables to identify the token.
    string public name = "CloudSlit Dao";
    string public symbol = "CSD";
    uint8 public decimals = 0;


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

    mapping(string => address) userBindings;

    function getWallet(string memory uid) external view returns(address){
        return userBindings[uid];
    }

    function bindWallet(string memory uid) external {
        require(userBindings[uid] == address(0));
        userBindings[uid] = msg.sender;
    }

    // Todo
    function changeWallet(string memory uid) external {
        userBindings[uid] = msg.sender;
    }
    
    //Deposit amount
    uint256 public fullnodeDepositAmount = 5000;
    uint256 public privoderDepositAmount = 1000;
    // // A mapping is a key/value map. Here we store each staked user.
    mapping(address => uint256) fullnodeDeposits;
    mapping(address => uint256) privateDeposits;

    // /**
    //  * 
    //  */
    function isDeposit(uint8 _type) external view returns (bool) {
        if(_type == 1){
            return fullnodeDeposits[msg.sender] != 0;
        } else if(_type == 2){
            return privateDeposits[msg.sender] != 0;
        }
        return false;
    }
    
    // /**
    //  * 
    //  */
    function stake(uint8 _type) external {
        if(_type == 1){
            require(fullnodeDeposits[msg.sender] == 0, "Already staked");
            require(balances[msg.sender] >= fullnodeDepositAmount, "Not enough tokens");
            balances[msg.sender] -= fullnodeDepositAmount;
            fullnodeDeposits[msg.sender] += fullnodeDepositAmount;
        }else if(_type == 2){
            require(privateDeposits[msg.sender] == 0, "Already staked");
            require(balances[msg.sender] >= privoderDepositAmount, "Not enough tokens");
            balances[msg.sender] -= privoderDepositAmount;
            privateDeposits[msg.sender] += privoderDepositAmount;
        }
    }
    // /**
    //  * 
    //  */
    function withdraw(uint8 _type) external {
        if(_type == 1){
            require(fullnodeDeposits[msg.sender] > 0);
            balances[msg.sender] += fullnodeDeposits[msg.sender];
            delete fullnodeDeposits[msg.sender];
        }else if(_type == 2){
            require(privateDeposits[msg.sender] > 0);
            balances[msg.sender] += privateDeposits[msg.sender];
            delete privateDeposits[msg.sender];
        }
    }
    
    struct Order {
        string name;
        uint startTime;
        uint endTime;
        uint withdrawTime;
        uint32 duration;
        uint256 price;
        bool used;
        bool withdraw;
        address payUser;
    }
    
    mapping(string=>Order) orders;
    mapping(address=>string[]) privoderOrders;

    function clientOrder(string memory _name, uint32 _duration, string memory _orderId, uint256 _price) external {
        require(orders[_orderId].used, "Already paid");
        require(balances[msg.sender] >= _price, "Not enough tokens");
        balances[msg.sender] -= _price;
        uint cooldownTime = 1 hours;
        orders[_orderId] = Order(_name, block.timestamp, block.timestamp + _duration * cooldownTime, 0, _duration, _price, true, false, msg.sender );
    }

    function checkOrder(string memory _orderId) public view returns(bool) {
        return (orders[_orderId].used);
    }

    function withdrawAllOrderTokens() external {
        require(privateDeposits[msg.sender] != 0);
        string[] memory _orders = privoderOrders[msg.sender];
        for (uint i=0; i < _orders.length; i++){
            if (!orders[_orders[i]].withdraw){
                
            }
        }
    }
    
    function withdrawOrderTokens(string memory _orderId) external {
        require(orders[_orderId].used);

    }
}