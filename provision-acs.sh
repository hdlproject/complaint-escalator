# Variables - Customize these
RESOURCE_GROUP="complaint-escalator-rg"
REGION="asiapacific"
LOCATION="southeastasia"
ACS_NAME="complaint-escalator-acs"
EMAIL_SERVICE_NAME="complaint-escalator-email-service"
DOMAIN_NAME="hdlproject.dev"  # For email (must already be owned by you)

# 1. Create a Resource Group
az group create \
  --name $RESOURCE_GROUP \
  --location $LOCATION

# 2. Create the Communication Services resource
az communication create \
  --name $ACS_NAME \
  --data-location $REGION \
  --resource-group $RESOURCE_GROUP \
  --location global

# 3. Get the connection string (needed for SDKs)
az communication list-key \
  --name $ACS_NAME \
  --resource-group $RESOURCE_GROUP

# 4. Create a DNS zone
az network dns zone create \
  --resource-group $RESOURCE_GROUP \
  --name $DOMAIN_NAME

# 5 Enable email service in the ACS
# Can only be done manually at this moment
# See az communication email domain initiate-verification for the possibility to automate this

# 6. Create an email domain
az communication email domain create \
  --name $DOMAIN_NAME \
  --resource-group $RESOURCE_GROUP \
  --email-service-name $EMAIL_SERVICE_NAME \
  --domain-management CustomerManaged

# 7. Wait for the domain to be verified
# If not, manually verify it in the UI

# 8. Show the email domain
az communication email domain show \
  --name $DOMAIN_NAME \
  --resource-group $RESOURCE_GROUP \
  --email-service-name $EMAIL_SERVICE_NAME

# 9. Connect the domain to the ACS and create a sender username
az communication email domain sender-username create \
  --resource-group $RESOURCE_GROUP \
  --email-service-name $EMAIL_SERVICE_NAME \
  --domain-name $DOMAIN_NAME \
  --name noreply \
  --username noreply
