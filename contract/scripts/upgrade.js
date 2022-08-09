async function main() {
    const proxyAddress = '0x45292fC8824A5cFF141B37E0e85c787390Ed76B7';

    const CloudSlitDaoV2 = await ethers.getContractFactory("CloudSlitDaoV2");
    console.log("Preparing upgrade...");

    await upgrades.upgradeProxy(proxyAddress, CloudSlitDaoV2);
    console.log("Upgraded Successfully");
}

main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });