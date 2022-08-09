async function main() {
    const proxyAddress = '0x3A35207918FEE0F59a32a1a36B58A758B4F222de';

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