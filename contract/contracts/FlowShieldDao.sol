// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract FlowShieldDao is ERC20, Ownable {

    struct userWallet {
        address user;
        uint8 status;
    }

    struct order {
        string name;
        uint startTime;
        uint endTime;
        uint withdrawDuration;
        uint32 duration;
        uint amount;
        bool used;
        bool withdraw;
        address payAddress;
        address privoderAddress;
    }

    mapping(string => userWallet) userWallets;

    //Initialize variables
    uint public _fullnodeDepositAmount;
    uint public _privoderDepositAmount;
    uint32 _durationUnit;
    // // A mapping is a key/value map. Here we store each staked user.
    mapping(address => uint) _fullnodeDeposits;
    mapping(address => uint) _privoderDeposits;

    mapping(string=>order) _orders;
    mapping(address=>string[]) _privoderOrders;

    constructor() ERC20("FlowShield Dao", "FSD") {
        _mint(msg.sender, 100000000 * 10 ** decimals());
        _fullnodeDepositAmount = 5000 * 10 ** decimals();
        _privoderDepositAmount = 1000 * 10 ** decimals();
        _durationUnit = 1 hours;
    }

    function getWallet(string memory uuid) external view returns(address, uint8){
        return (userWallets[uuid].user, userWallets[uuid].status);
    }

    function bindWallet(string memory uuid) external {
        require(userWallets[uuid].user == address(0));
        if (_fullnodeDeposits[msg.sender] == 0) {
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
        require(_fullnodeDeposits[msg.sender] > 0);
        require(userWallets[uuid].status == 1);
        userWallets[uuid].status = 2;
    }

    function changeWallet(string memory uuid, address newWallet) external {
        require(newWallet != address(0));
        if (userWallets[uuid].status == 1){
            userWallets[uuid].user = newWallet;
        }else{
            require(userWallets[uuid].user == msg.sender);
            userWallets[uuid].user = newWallet;
        }
    }

    function getUserInfo(string memory uuid) external view returns(bool, bool){
        if(userWallets[uuid].status == 2){
            return ((_fullnodeDeposits[userWallets[uuid].user] > 0), (_privoderDeposits[userWallets[uuid].user] > 0));
        }else{
            return (false, false);
        }
    }
    // /**
    //  *
    //  */
    function isDeposit(uint8 _type) external view returns (bool) {
        if(_type == 1){
            return _fullnodeDeposits[msg.sender] != 0;
        } else if(_type == 2){
            return _privoderDeposits[msg.sender] != 0;
        }
        return false;
    }

    // /**
    //  *
    //  */
    function getDeposit(address walletAddress) external view returns (uint, uint) {
        return (_fullnodeDeposits[walletAddress], _privoderDeposits[walletAddress]);
    }

    function stakeAmount(uint8 _type, address walletAddress, uint amount) external {
        if(_type == 1){
            require(balanceOf(msg.sender) >= amount, "Not enough CSD");
            transfer(address(this), amount);
            _fullnodeDeposits[walletAddress] += amount;
        }else if(_type == 2){
            require(balanceOf(msg.sender) >= amount, "Not enough CSD");
            transfer(address(this), amount);
            _privoderDeposits[walletAddress] += amount;
        }
    }
    // /**
    //  *
    //  */
    function stake(uint8 _type) external {
        if(_type == 1){
            require(_fullnodeDeposits[msg.sender] == 0, "Already staked");
            require(balanceOf(msg.sender) >= _fullnodeDepositAmount, "Not enough CSD");
            transfer(address(this), _fullnodeDepositAmount);
            _fullnodeDeposits[msg.sender] += _fullnodeDepositAmount;
        }else if(_type == 2){
            require(_privoderDeposits[msg.sender] == 0, "Already staked");
            require(balanceOf(msg.sender) >= _privoderDepositAmount, "Not enough CSD");
            transfer(address(this), _privoderDepositAmount);
            _privoderDeposits[msg.sender] += _privoderDepositAmount;
        }
    }
    // /**
    //  *
    //  */
    function withdraw(uint8 _type) external {
        if(_type == 1){
            require(_fullnodeDeposits[msg.sender] > 0);
            transferFrom(address(this), msg.sender, _fullnodeDeposits[msg.sender]);
            delete _fullnodeDeposits[msg.sender];
        }else if(_type == 2){
            require(_privoderDeposits[msg.sender] > 0);
            transferFrom(address(this), msg.sender, _privoderDeposits[msg.sender]);
            delete _privoderDeposits[msg.sender];
        }
    }

    function checkOrder(string memory orderId) public view returns(bool) {
        return (_orders[orderId].used);
    }

    function getOrdersInfo(string memory orderId) public view returns(order memory){
        return (_orders[orderId]);
    }

    function clientOrder(string memory name, uint32 duration, string memory orderId, uint amount, address to) external {
        require(!_orders[orderId].used, "Already paid");
        require(balanceOf(msg.sender) >= amount, "Not enough CSD");
        transfer(address(this), amount);
        _orders[orderId] = order(name, block.timestamp, block.timestamp + duration * _durationUnit, 0, duration, amount, true, false, msg.sender , to);
        _privoderOrders[to].push(orderId);
    }

    function getPrivoderOrders(address from) public view returns(string[] memory ){
        return _privoderOrders[from];
    }

    function getAllOrderTokens() external view returns(uint){
        if (_privoderDeposits[msg.sender] == 0){
            return 0;
        }
        string[] memory orders = _privoderOrders[msg.sender];
        uint amount = 0;
        for (uint i=0; i < orders.length; i++){
            if (!_orders[orders[i]].withdraw){
                if(block.timestamp >= _orders[orders[i]].endTime){
                    uint duration = _orders[orders[i]].duration  - _orders[orders[i]].withdrawDuration;
                    amount += (_orders[orders[i]].amount / _orders[orders[i]].duration) * duration;
                }else{
                    uint duration = (block.timestamp - _orders[orders[i]].startTime) / _durationUnit  - _orders[orders[i]].withdrawDuration;
                    amount += (_orders[orders[i]].amount / _orders[orders[i]].duration) * duration;
                }
            }
        }
        return (amount);
    }

    function withdrawAllOrderTokens() external {
        require(_privoderDeposits[msg.sender] != 0, 'Not deposits');
        string[] memory orders = _privoderOrders[msg.sender];
        uint amount = 0;
        for (uint i=0; i < orders.length; i++){
            if (!_orders[orders[i]].withdraw){
                if(block.timestamp >= _orders[orders[i]].endTime){
                    _orders[orders[i]].withdraw = true;
                    uint duration = _orders[orders[i]].duration  - _orders[orders[i]].withdrawDuration;
                    amount += (_orders[orders[i]].amount / _orders[orders[i]].duration) * duration;
                }else{
                    uint duration = (block.timestamp - _orders[orders[i]].startTime) / _durationUnit  - _orders[orders[i]].withdrawDuration;
                    _orders[orders[i]].withdrawDuration += duration;
                    amount += (_orders[orders[i]].amount / _orders[orders[i]].duration) * duration;
                }
            }
        }
        transferFrom(address(this), msg.sender, amount);
    }

    function withdrawOrderTokens(string memory orderId) external {
        require(!_orders[orderId].withdraw, 'The order has been withdrawn');
        require(_orders[orderId].privoderAddress == msg.sender, 'Please confirm the wallet address, Can not withdraw');
        if(block.timestamp >= _orders[orderId].endTime){
            _orders[orderId].withdraw = true;
            uint duration = _orders[orderId].duration  - _orders[orderId].withdrawDuration;
            uint amount = (_orders[orderId].amount / _orders[orderId].duration) * duration;
            transferFrom(address(this), msg.sender, amount);
        }else{
            uint duration = (block.timestamp - _orders[orderId].startTime) / _durationUnit  - _orders[orderId].withdrawDuration;
            _orders[orderId].withdrawDuration += duration;
            uint amount = (_orders[orderId].amount / _orders[orderId].duration) * duration;
            transferFrom(address(this), msg.sender, amount);
        }
    }
}