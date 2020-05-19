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
func GetVMResourcegroup(vm *compute.VirtualMachine)(string,error){
	   var resourceGroup string
	   if vm.ID !=nil{
		   s := strings.Split(*vm.ID,"/")
		   resourceGroup = s[4]
		   return resourceGroup,nil
	   }else{
			 return resourceGroup,errors.New("No resourceGroup")
	   }	 
}
//Returns the virtual machine's name
func GetVMname(vm *compute.VirtualMachine)(string,error){
	var Name string
	if vm.ID !=nil{
		s := strings.Split(*vm.ID,"/")
		Name = s[8]
		return Name,nil
	}else{
		  return Name,errors.New("No vm name")
	}	 
}
//Returns the subscription ID
func GetVMSubscription(vm *compute.VirtualMachine)(string,error){
	var subscription string
	if vm.ID !=nil{
		s := strings.Split(*vm.ID,"/")
		subscription = s[2]
		return subscription,nil
	}else{
		  return subscription,errors.New("No subscription")
	}	 
}
//Returns the tags related to the virtual machine
func GetVMTags(vm *compute.VirtualMachine)(map[string]*string,error){
	var tags map[string]*string
	if vm.Tags !=nil{
	
		tags = vm.Tags
		return tags,nil
	}else{
		  return tags,errors.New("no tags present for the vm")
	}	 
}
//Returns the Location 
func GetVMLocation(vm *compute.VirtualMachine)(string,error){
	var location string
	if vm.Location !=nil{
		location =  *vm.Location
		return location,nil
	}else{
		  return location,errors.New("no location assigned to the vm")
	}	 
}

func GetVMSize(vm *compute.VirtualMachine)(compute.VirtualMachineSizeTypes){
	
		Vmsize:= vm.VirtualMachineProperties.HardwareProfile.VMSize
		return Vmsize
	 
}
//Returns the OStype used in th virtual machine
func GetVMOsType(vm *compute.VirtualMachine)(compute.OperatingSystemTypes){
	
		VmOS:= vm.VirtualMachineProperties.StorageProfile.OsDisk.OsType
		return VmOS
}

func GetVMadminusername(vm *compute.VirtualMachine)(string){
	
	VMadminusername:= *vm.VirtualMachineProperties.OsProfile.AdminUsername
	return VMadminusername
}

func GetVmnetworkinterface(vm *compute.VirtualMachine)(string,error){
	networkinterface:=*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces
	netinterface:=*networkinterface[0].ID
	ID := strings.Split(netinterface,"/")
	netwinterface := ID[8]
	return netwinterface,nil
}
//Returns the publicIPname 
func GetPublicIPname(subscriptionID string,resourceGroup string,networkinterface string)(string,error){
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}
	PublicIPname:=""
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	vmInterface := network.NewInterfacesClient(subscriptionID)
	vmInterface.Authorizer = authorizer
	interfaces,err:= vmInterface.Get(ctx,resourceGroup,networkinterface,"")
	if err != nil {
		panic(err)
	}
	interfaceinfo :=*interfaces.InterfacePropertiesFormat.IPConfigurations
	interfID := *interfaceinfo[0].InterfaceIPConfigurationPropertiesFormat
	if interfID.PublicIPAddress!=nil&&interfID.PublicIPAddress.ID!=nil {
	    ID:=strings.Split(*interfID.PublicIPAddress.ID,"/")
		return ID[8],nil
	}else{
		return PublicIPname , errors.New("Vm has no PublicIPname")
	}
}
//Returns the publicIPAddress of the virtual machine
func GetPublicIpaddress(subscriptionID string,resourceGroup string,PublicIPname string)(string,error){
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}
	PublicIPaddress:=""
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	vmPublicIP:=network.NewPublicIPAddressesClient(subscriptionID)
	vmPublicIP.Authorizer = authorizer
	VmIP,err := vmPublicIP.Get(ctx,resourceGroup,PublicIPname,"")
	if err != nil {
		panic(err)
	}
	if VmIP.PublicIPAddressPropertiesFormat!=nil && VmIP.PublicIPAddressPropertiesFormat.IPAddress!=nil{
		PublicIPAddress:=*VmIP.PublicIPAddressPropertiesFormat.IPAddress
	    return PublicIPAddress,nil
	}else{
					
		return PublicIPaddress,errors.New("Vm has no PublicIPaddress")

	}
}
//Returns the virtual network and subnet
func GetSubnetandvirtualnetwork(subscriptionID string,resourceGroup string,networkinterface string)(string,error){
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}
	subnet:=""
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	vmInterface := network.NewInterfacesClient(subscriptionID)
	vmInterface.Authorizer = authorizer
	interfaces,err:= vmInterface.Get(ctx,resourceGroup,networkinterface,"")
	if err != nil {
		panic(err)
	}
	interfaceinfo :=*interfaces.InterfacePropertiesFormat.IPConfigurations
	interfID := *interfaceinfo[0].InterfaceIPConfigurationPropertiesFormat
	if interfID.Subnet!=nil {
		ID := strings.Split(*interfID.Subnet.ID,"/")
		virtualnetwork:= ID[8]+"/"+ID[10]
		return virtualnetwork , nil
	}else{
		return subnet,errors.New("Vm has no virtual network and subnet")
	}
}
