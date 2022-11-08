// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// This is the main building block for smart contracts.
contract ERC20 {
    string public name = "CloudSlit Dao";
    string public symbol = "CSD";
    uint256 public totalSupply;
    uint8 public decimals = 18;

    mapping(address => uint256) private balances;
    mapping(address => mapping(address => uint256)) private allowances;


    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    function _mint(address account, uint256 amount) internal virtual {
        require(account != address(0), "ERC20: mint to the zero address");

        totalSupply += amount;
        // Overflow not possible: balance + amount is at most totalSupply + amount, which is checked above.
        balances[account] += amount;
    }

    function transfer(address to, uint256 amount) public virtual returns (bool){
        // Check if the transaction sender has enough tokens.
        // If `require`'s first argument evaluates to `false` then the
        // transaction will revert.
        require(balances[msg.sender] >= amount, "Not enough tokens");

        // Transfer the amount.
        _transfer(msg.sender, to, amount);
        return true;
    }
    /**
         * @dev Moves `amount` of tokens from `from` to `to`.
     *
     * This internal function is equivalent to {transfer}, and can be used to
     * e.g. implement automatic token fees, slashing mechanisms, etc.
     *
     * Emits a {Transfer} event.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `from` must have a balance of at least `amount`.
     */
    function _transfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual {
        require(from != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");

        uint256 fromBalance = balances[from];
        require(fromBalance >= amount, "ERC20: transfer amount exceeds balance");
        balances[from] = fromBalance - amount;
        // Overflow not possible: the sum of all balances is capped by totalSupply, and the sum is preserved by
        // decrementing then incrementing.
        balances[to] += amount;
    }

    function approve(address spender, uint256 value) public returns (bool){
        allowances[msg.sender][spender] = value;
        return true;
    }

    function allowance(address owner, address spender) public view returns (uint256){
        return allowances[owner][spender];
    }

    function balanceOf(address account) public view returns (uint256) {
        return balances[account];
    }
}