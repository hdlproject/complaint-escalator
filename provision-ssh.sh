RESOURCE_GROUP="complaint-escalator-rg"
VM_NAME="complaint-escalator-vm"
VM_USERNAME="hdl"
VM_PASSWORD="Hdl@12345678"
LOCATION="southeastasia"

# Create the VM
az vm create \
  --resource-group $RESOURCE_GROUP \
  --name $VM_NAME \
  --image Ubuntu2404 \
  --admin-username $VM_USERNAME \
  --admin-password $VM_PASSWORD \
  --authentication-type ssh \
  --generate-ssh-keys \
  --location $LOCATION \
  --size Standard_D2d_v4

# # Update the VM username and password
# az vm user update \
#   --resource-group $RESOURCE_GROUP \
#   --name $VM_NAME \
#   --username $VM_USERNAME \
#   --password $VM_PASSWORD

# # Update the VM SSH key
# az vm user update \
#   --resource-group $RESOURCE_GROUP \
#   --name $VM_NAME \
#   --username $VM_USERNAME \
#   --ssh-key-value ~/.ssh/id_ed25519.pub

az ssh vm --resource-group $RESOURCE_GROUP --name $VM_NAME --local-user $VM_USERNAME
