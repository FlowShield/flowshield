async function main() {
    const CloudSlitDao = await ethers.getContractFactory("CloudSlitDao");
    console.log("Deploying CloudSlitDao...");

    const contract = await upgrades.deployProxy(CloudSlitDao, { initializer: 'initialize', kind: "transparent",});
    await contract.deployed();
    console.log("CloudSlitDao deployed to:", contract.address);
}

main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });