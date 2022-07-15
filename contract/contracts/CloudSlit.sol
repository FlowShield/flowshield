// Solidity files have to start with this pragma.
// It will be used by the Solidity compiler to validate its version.
pragma solidity ^0.8.15;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// This is the main building block for smart contracts.
contract CloudSlit is ERC20 {

    constructor(uint256 initialSupply) ERC20("CloudSlit Dao", "CSD") public {
        _mint(msg.sender, initialSupply * 10 ** decimals());
    }

    struct userWallet {
        address user;
        uint8 status;
    }
    mapping(string => userWallet) userWallets;

    function getWallet(string memory uuid) external view returns(address, uint8){
        return (userWallets[uuid].user, userWallets[uuid].status);
    }

    function bindWallet(string memory uuid) external {
        require(userWallets[uuid].user == address(0));
        if (fullnodeDeposits[msg.sender] == 0) {
            userWallets[uuid] = userWallet(msg.sender, 1);
        }else{
            userWallets[uuid] = userWallet(msg.sender, 2);
        }
    }

    function unbindWallet(string memory uuid) external {
        require(userWallets[uuid].user == msg.sender);
        delete userWallets[uuid];
    }

    function verifyWallet(string memory uuid) external {
        require(fullnodeDeposits[msg.sender] > 0);
        require(userWallets[uuid].status == 1);
        userWallets[uuid].status = 2;
    }

    // Todo
    function changeWallet(string memory uuid, address newWallet) external {
        require(newWallet != address(0));
        if (userWallets[uuid].status == 1){
            userWallets[uuid].user = newWallet;
        }else{
            require(userWallets[uuid].user == msg.sender);
            userWallets[uuid].user = newWallet;
        }
    }
    
    //Deposit amount
    uint256 public fullnodeDepositAmount = 5000;
    uint256 public privoderDepositAmount = 1000;
    // // A mapping is a key/value map. Here we store each staked user.
    mapping(address => uint256) fullnodeDeposits;
    mapping(address => uint256) privateDeposits;

    function getUserInfo(string memory uuid) external view returns(bool, bool){
        if(userWallets[uuid].status == 2){
            return ((fullnodeDeposits[userWallets[uuid].user] > 0), (privateDeposits[userWallets[uuid].user] > 0));
        }else{
            return (false, false);
        }
    }
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
            require(balanceOf(msg.sender) >= fullnodeDepositAmount, "Not enough CSD");
            _transfer(msg.sender, address(this), fullnodeDepositAmount);
            fullnodeDeposits[msg.sender] += fullnodeDepositAmount;
        }else if(_type == 2){
            require(privateDeposits[msg.sender] == 0, "Already staked");
            require(balanceOf(msg.sender) >= privoderDepositAmount, "Not enough CSD");
            _transfer(msg.sender, address(this), privoderDepositAmount);
            privateDeposits[msg.sender] += privoderDepositAmount;
        }
    }
    // /**
    //  * 
    //  */
    function withdraw(uint8 _type) external {
        if(_type == 1){
            require(fullnodeDeposits[msg.sender] > 0);
            _transfer(address(this), msg.sender, fullnodeDeposits[msg.sender]);
            delete fullnodeDeposits[msg.sender];
        }else if(_type == 2){
            require(privateDeposits[msg.sender] > 0);
            _transfer(address(this), msg.sender, privateDeposits[msg.sender]);
            delete privateDeposits[msg.sender];
        }
    }
    
    struct Order {
        string name;
        uint startTime;
        uint endTime;
        uint withdrawDuration;
        uint32 duration;
        uint price;
        bool used;
        bool withdraw;
        address payAddress;
        address privateAddress;
    }

    function getOrdersInfo(string memory _orderId) public view returns(Order memory){
        return (orders[_orderId]);
    }

    mapping(string=>Order) orders;
    mapping(address=>string[]) privoderOrders;
    uint32 durationUnit = 1 hours;

    function clientOrder(string memory _name, uint32 _duration, string memory _orderId, uint256 _price, address _to) external {
        require(!orders[_orderId].used, "Already paid");
        require(balanceOf(msg.sender) >= _price, "Not enough CSD");
        _transfer(msg.sender, address(this), _price);
        orders[_orderId] = Order(_name, block.timestamp, block.timestamp + _duration * durationUnit, 0, _duration, _price, true, false, msg.sender , _to);
        privoderOrders[_to].push(_orderId);
    }

    function checkOrder(string memory _orderId) public view returns(bool) {
        return (orders[_orderId].used);
    }

    function getPrivoderOrders(address from) public view returns(string[] memory ){
        return privoderOrders[from];
    }

    function getAllOrderTokens() external view returns(uint){
        if (privateDeposits[msg.sender] != 0){
            return 0;
        }
        string[] memory _orders = privoderOrders[msg.sender];
        uint price = 0;
        for (uint i=0; i < _orders.length; i++){
            if (!orders[_orders[i]].withdraw){
                if(block.timestamp >= orders[_orders[i]].endTime){
                    uint duration = orders[_orders[i]].duration  - orders[_orders[i]].withdrawDuration;
                    price += (orders[_orders[i]].price / orders[_orders[i]].duration) * duration;
                }else{
                    uint duration = (block.timestamp - orders[_orders[i]].startTime) / durationUnit  - orders[_orders[i]].withdrawDuration;
                    price += (orders[_orders[i]].price / orders[_orders[i]].duration) * duration;
                }
            }
        }
        return (price);
    }

    function withdrawAllOrderTokens() external {
        require(privateDeposits[msg.sender] != 0, 'Not deposits');
        string[] memory _orders = privoderOrders[msg.sender];
        uint price = 0;
        for (uint i=0; i < _orders.length; i++){
            if (!orders[_orders[i]].withdraw){
                if(block.timestamp >= orders[_orders[i]].endTime){
                    orders[_orders[i]].withdraw = true;
                    uint duration = orders[_orders[i]].duration  - orders[_orders[i]].withdrawDuration;
                    price += (orders[_orders[i]].price / orders[_orders[i]].duration) * duration;
                }else{
                    uint duration = (block.timestamp - orders[_orders[i]].startTime) / durationUnit  - orders[_orders[i]].withdrawDuration;
                    orders[_orders[i]].withdrawDuration += duration;
                    price += (orders[_orders[i]].price / orders[_orders[i]].duration) * duration;
                }
            }
        }
        _transfer(address(this), msg.sender, price);
    }

    function withdrawOrderTokens(string memory _orderId) external {
        require(!orders[_orderId].withdraw, 'The order has been withdrawn');
        require(orders[_orderId].privateAddress == msg.sender, 'Please confirm the wallet address, Can not withdraw');
        if(block.timestamp >= orders[_orderId].endTime){
            orders[_orderId].withdraw = true;
            uint duration = orders[_orderId].duration  - orders[_orderId].withdrawDuration;
            uint price = (orders[_orderId].price / orders[_orderId].duration) * duration;
            _transfer(address(this), msg.sender, price);
        }else{
            uint duration = (block.timestamp - orders[_orderId].startTime) / durationUnit  - orders[_orderId].withdrawDuration;
            orders[_orderId].withdrawDuration += duration;
            uint price = (orders[_orderId].price / orders[_orderId].duration) * duration;
            _transfer(address(this), msg.sender, price);
        }
    }
}