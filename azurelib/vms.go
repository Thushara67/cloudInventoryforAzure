package azurelib

import (
	"context"
	"time"
	"strings"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

//A struct that contains all the necessary clients
type Clients struct {
	//Network Interface Client
	VmInterface  network.InterfacesClient
	//Public IP Addresses Client
	VmPublicIP network.PublicIPAddressesClient
	//Virtual Machine Client
	VmClient compute.VirtualMachinesClient
}

//Returns a New Client 
//Parameters - subscriptionID : Subscription ID for Azure
func GetNewClients (subscriptionID string) Clients {
	VmInterface := network.NewInterfacesClient(subscriptionID)
	VmPublicIP:=network.NewPublicIPAddressesClient(subscriptionID)
	VmClient := compute.NewVirtualMachinesClient(subscriptionID)
	
	c := Clients{VmInterface, VmPublicIP, VmClient}
	return c
}

//Authorizes all the clients
func AuthorizeClients (c Clients) Clients{
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}
	c.VmClient.Authorizer = authorizer
	c.VmPublicIP.Authorizer = authorizer
	c.VmInterface.Authorizer = authorizer	
	return c
}

//Get Private IP Address of a Virtual Machine
func GetPrivateIP (vmInterface network.InterfacesClient, ctx context.Context, 
	resourceGroup string, networkInterface string, expand string) (PrivateIPAddress string, 
	IPConfiguration string, err error) {
	interfaces,err:= vmInterface.Get(ctx,resourceGroup,networkInterface,expand)
	if err != nil {
		return 
	}
	interfaceinfo :=*interfaces.InterfacePropertiesFormat.IPConfigurations
	interfID := *interfaceinfo[0].InterfaceIPConfigurationPropertiesFormat
	//fmt.Println("IP configuration :",*interfaceinfo[0].Name)
	IPConfiguration = *interfaceinfo[0].Name
	if interfID.PrivateIPAddress!=nil {
	//fmt.Println("PrivateIpaddress :",*interfID.PrivateIPAddress)
		PrivateIPAddress = *interfID.PrivateIPAddress
	}
	return
}

//Get Public IP Address ID (PublicIPName)
func GetPublicIPAddressID (vmInterface network.InterfacesClient, 
	ctx context.Context, resourceGroup string, networkInterface string, 
	expand string) (PublicIPAddressID string, err error) {
	interfaces,err:= vmInterface.Get(ctx,resourceGroup,networkInterface,expand)
	if err != nil {
		return 
	}
	interfaceinfo :=*interfaces.InterfacePropertiesFormat.IPConfigurations
	interfID := *interfaceinfo[0].InterfaceIPConfigurationPropertiesFormat
	
	if interfID.PublicIPAddress!=nil&&interfID.PublicIPAddress.ID!=nil {
		ID:=strings.Split(*interfID.PublicIPAddress.ID,"/")
		//fmt.Println("PublicIPAddress ID : ",ID[8])
		PublicIPAddressID = ID[8]		
	}else{
		err = errors.New("Vm has no publicIPname")
	}
	return
}


func GetallVMS(subscriptionID string)([]*compute.VirtualMachine,error){

    authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}

	vmClient := compute.NewVirtualMachinesClient(subscriptionID)
	vmClient.Authorizer = authorizer
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	results, err := vmClient.ListAllComplete(ctx)
	if err != nil {
		panic(err)
	}
	var Vmlist []*compute.VirtualMachine
    for results.NotDone(){
		vm := results.Value()
		Vmlist = append(Vmlist,&vm)
        if err := results.Next(); err != nil {
			panic(err)
		}
		
	}

	return Vmlist,nil

}
//Returns resourcegroup to which the virtual machine belongs to
func GetVMResourcegroup(vm *compute.VirtualMachine)(resourceGroup string,err error){
	   
	if vm.ID !=nil{
		s := strings.Split(*vm.ID,"/")
		resourceGroup = s[4]
		err = nil
	}else{
		  err = errors.New("No resourceGroup")
	}	 
	return
}
//Returns the virtual machine's name
func GetVMname(vm *compute.VirtualMachine)(Name string,err error){
 
 if vm.ID !=nil{
	 s := strings.Split(*vm.ID,"/")
	 Name = s[8]
	 err =nil
 }else{
	  err = errors.New("No vm name")
 }	 
 return
}
//Returns the subscription ID
func GetVMSubscription(vm *compute.VirtualMachine)(subscriptionID string,err error){
 
 if vm.ID !=nil{
	 s := strings.Split(*vm.ID,"/")
	 subscriptionID = s[2]
	 err = nil
 }else{
	   err = errors.New("No subscription")
 }	 
 return
}
//Returns the tags related to the virtual machine
func GetVMTags(vm *compute.VirtualMachine)(tags map[string]*string,err error){
 if vm.Tags !=nil{
 
	 tags = vm.Tags
	 err = nil
 }else{
	   err = errors.New("no tags present for the vm")
 }	
 return  
}
//Returns the Location 
func GetVMLocation(vm *compute.VirtualMachine)(location string,err error){
 
 if vm.Location !=nil{
	 location =  *vm.Location
	 err = nil
 }else{
	   err = errors.New("no location assigned to the vm")
 }	
 return  
}

func GetVMSize(vm *compute.VirtualMachine)(Vmsize compute.VirtualMachineSizeTypes){
 
	 Vmsize = vm.VirtualMachineProperties.HardwareProfile.VMSize
	 return 
  
}
//Returns the OStype used in th virtual machine
func GetVMOsType(vm *compute.VirtualMachine)(VmOS compute.OperatingSystemTypes){
 
	 VmOS = vm.VirtualMachineProperties.StorageProfile.OsDisk.OsType
	 return 
}

func GetVMadminusername(vm *compute.VirtualMachine)(VMadminusername string,err error){
 if vm.VirtualMachineProperties.OsProfile.AdminUsername!=nil{
	VMadminusername = *vm.VirtualMachineProperties.OsProfile.AdminUsername
	err = nil
 }else{
	 err = errors.New("Vm has no admin user name")
 }
 return
}

func GetVmnetworkinterface(vm *compute.VirtualMachine)(networkInterface string,err error){
 if vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces!=nil{
	 networkinterface:=*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces
	 netinterface:=*networkinterface[0].ID
	 ID := strings.Split(netinterface,"/")
	 networkInterface = ID[8]
	 err = nil
 }else{
	 err = errors.New("Vm has no network interface")
 }
 return
 
}

//Returns the PublicIPAddress of the virtual machine
func GetPublicIPAddress(vmPublicIP network.PublicIPAddressesClient, ctx context.Context,
 resourceGroup string, PublicIPname string, expand string) (PublicIPAddress string, err error) {
 VmIP,err := vmPublicIP.Get(ctx, resourceGroup, PublicIPname, expand)
 if err != nil {
	 return
 }
 if VmIP.PublicIPAddressPropertiesFormat!=nil && VmIP.PublicIPAddressPropertiesFormat.IPAddress!=nil{
	 PublicIPAddress = *VmIP.PublicIPAddressPropertiesFormat.IPAddress
	 
 }else{
	 err = errors.New("Vm has no publicIPAddress")
 }
 return 
 
}

//Returns the virtual network and subnet
func GetSubnetandvirtualnetwork(vmInterface network.InterfacesClient, 
 ctx context.Context,resourceGroup string,networkinterface string,expand string)(virtualnetworkandsubnet string,err error){
 interfaces,err:= vmInterface.Get(ctx,resourceGroup,networkinterface,expand)
 if err != nil {
	 return
 }
 interfaceinfo :=*interfaces.InterfacePropertiesFormat.IPConfigurations
 interfID := *interfaceinfo[0].InterfaceIPConfigurationPropertiesFormat
 if interfID.Subnet!=nil {
	 ID := strings.Split(*interfID.Subnet.ID,"/")
	 virtualnetworkandsubnet =  ID[8]+"/"+ID[10]
 }else{
	 err = errors.New("Vm has no virtual network and subnet")
 }
 return
}

func GetDNS(vmPublicIP network.PublicIPAddressesClient, ctx context.Context,
 resourceGroup string, PublicIPname string, expand string)(Fqdn string,err error) {	
 VmIP,err := vmPublicIP.Get(ctx,resourceGroup, PublicIPname,expand)
 if err != nil {
	 return
 }
 if VmIP.PublicIPAddressPropertiesFormat!=nil && VmIP.PublicIPAddressPropertiesFormat.DNSSettings!=nil{
	 Fqdn = *VmIP.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn
 }else{
	 err = errors.New("DNS is not configured")
 }
 return
}
