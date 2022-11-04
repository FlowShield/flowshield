// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract CloudSlitDao is Initializable, ERC20Upgradeable, ERC20BurnableUpgradeable, OwnableUpgradeable {

    struct userWallet {
        address user;
        uint8 status;
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
        address privoderAddress;
    }

    mapping(string => userWallet) userWallets;

    //Initialize variables
    uint256 public _fullnodeDepositAmount;
    uint256 public _privoderDepositAmount;
    uint32 _durationUnit;
    // // A mapping is a key/value map. Here we store each staked user.
    mapping(address => uint256) _fullnodeDeposits;
    mapping(address => uint256) _privoderDeposits;

    mapping(string=>Order) _orders;
    mapping(address=>string[]) _privoderOrders;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() initializer public {
        __ERC20_init("CloudSlit Dao", "CSD");
        __ERC20Burnable_init();
        __Ownable_init();

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
    function stake(uint8 _type) external {
        if(_type == 1){
            require(_fullnodeDeposits[msg.sender] == 0, "Already staked");
            require(balanceOf(msg.sender) >= _fullnodeDepositAmount, "Not enough CSD");
            _transfer(msg.sender, address(this), _fullnodeDepositAmount);
            _fullnodeDeposits[msg.sender] += _fullnodeDepositAmount;
        }else if(_type == 2){
            require(_privoderDeposits[msg.sender] == 0, "Already staked");
            require(balanceOf(msg.sender) >= _privoderDepositAmount, "Not enough CSD");
            _transfer(msg.sender, address(this), _privoderDepositAmount);
            _privoderDeposits[msg.sender] += _privoderDepositAmount;
        }
    }
    // /**
    //  *
    //  */
    function withdraw(uint8 _type) external {
        if(_type == 1){
            require(_fullnodeDeposits[msg.sender] > 0);
            _transfer(address(this), msg.sender, _fullnodeDeposits[msg.sender]);
            delete _fullnodeDeposits[msg.sender];
        }else if(_type == 2){
            require(_privoderDeposits[msg.sender] > 0);
            _transfer(address(this), msg.sender, _privoderDeposits[msg.sender]);
            delete _privoderDeposits[msg.sender];
        }
    }

    function get_ordersInfo(string memory orderId) public view returns(Order memory){
        return (_orders[orderId]);
    }

    function clientOrder(string memory name, uint32 duration, string memory orderId, uint256 price, address to) external {
        require(!_orders[orderId].used, "Already paid");
        require(balanceOf(msg.sender) >= price, "Not enough CSD");
        _transfer(msg.sender, address(this), price);
        _orders[orderId] = Order(name, block.timestamp, block.timestamp + duration * _durationUnit, 0, duration, price, true, false, msg.sender , to);
        _privoderOrders[to].push(orderId);
    }

    function checkOrder(string memory orderId) public view returns(bool) {
        return (_orders[orderId].used);
    }

    function get_privoderOrders(address from) public view returns(string[] memory ){
        return _privoderOrders[from];
    }

    function getAllOrderTokens() external view returns(uint){
        if (_privoderDeposits[msg.sender] == 0){
            return 0;
        }
        string[] memory orders = _privoderOrders[msg.sender];
        uint price = 0;
        for (uint i=0; i < orders.length; i++){
            if (!_orders[orders[i]].withdraw){
                if(block.timestamp >= _orders[orders[i]].endTime){
                    uint duration = _orders[orders[i]].duration  - _orders[orders[i]].withdrawDuration;
                    price += (_orders[orders[i]].price / _orders[orders[i]].duration) * duration;
                }else{
                    uint duration = (block.timestamp - _orders[orders[i]].startTime) / _durationUnit  - _orders[orders[i]].withdrawDuration;
                    price += (_orders[orders[i]].price / _orders[orders[i]].duration) * duration;
                }
            }
        }
        return (price);
    }

    function withdrawAllOrderTokens() external {
        require(_privoderDeposits[msg.sender] != 0, 'Not deposits');
        string[] memory orders = _privoderOrders[msg.sender];
        uint price = 0;
        for (uint i=0; i < orders.length; i++){
            if (!_orders[orders[i]].withdraw){
                if(block.timestamp >= _orders[orders[i]].endTime){
                    _orders[orders[i]].withdraw = true;
                    uint duration = _orders[orders[i]].duration  - _orders[orders[i]].withdrawDuration;
                    price += (_orders[orders[i]].price / _orders[orders[i]].duration) * duration;
                }else{
                    uint duration = (block.timestamp - _orders[orders[i]].startTime) / _durationUnit  - _orders[orders[i]].withdrawDuration;
                    _orders[orders[i]].withdrawDuration += duration;
                    price += (_orders[orders[i]].price / _orders[orders[i]].duration) * duration;
                }
            }
        }
        _transfer(address(this), msg.sender, price);
    }

    function withdrawOrderTokens(string memory orderId) external {
        require(!_orders[orderId].withdraw, 'The order has been withdrawn');
        require(_orders[orderId].privoderAddress == msg.sender, 'Please confirm the wallet address, Can not withdraw');
        if(block.timestamp >= _orders[orderId].endTime){
            _orders[orderId].withdraw = true;
            uint duration = _orders[orderId].duration  - _orders[orderId].withdrawDuration;
            uint price = (_orders[orderId].price / _orders[orderId].duration) * duration;
            _transfer(address(this), msg.sender, price);
        }else{
            uint duration = (block.timestamp - _orders[orderId].startTime) / _durationUnit  - _orders[orderId].withdrawDuration;
            _orders[orderId].withdrawDuration += duration;
            uint price = (_orders[orderId].price / _orders[orderId].duration) * duration;
            _transfer(address(this), msg.sender, price);
        }
    }
}