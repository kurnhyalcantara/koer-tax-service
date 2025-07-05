# Qcash Base Protobuf

This repository contains the protobuf files for the Qcash base grpc protobuf files.

## How to use

1. Clone this repository
2. Create a new folder with the name of the service
3. Copy your proto files to the folder
4. Run the following command to verify the proto files compile correctly

```shell
./verify.sh module-name
```

5. If the proto files compile correctly, you can now use the generated go files in your destination service.

Notes :

If import path is not correct, you can adjust the import path with prefix `protos/module-name` in the proto files.

## How to integrate with other services

### Not integrated yet

1. Run the following command to add this repository as a submodule in your service:

```bash
git submodule add https://bitbucket.bri.co.id/scm/bricams-addons/qcash-base-protobuf.git proto
```

2. Copy `scripts/generate.sh` to the root of your service
3. Adjust the `generate.sh` file to match the proto files in your service
4. Run the following command to generate the go files:

```bash
./generate.sh module-name1 module-name2
```

5. Update the import path in the proto files to match the new import path

### Already integrated, but first time cloning the repository

Run the following command to clone the submodule:

```bash
git submodule init
git submodule update
```

### Already integrated, but updating the proto files

Run the following command to update the submodule:

```bash
git submodule update --remote --merge
```

Example : See [RnD Service](https://bitbucket.bri.co.id/projects/BRICAMS-ADDONS/repos/qcash-rnd-service)

**Created and Maintained by @Roxanne Squad**