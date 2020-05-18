package main

import (
	"context"
	"fmt"
	"time"
	"strings"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const subscriptionID = "282160c0-3c83-43f1-bff1-9356b1678ffb"

func main() {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}

	vmInterface := network.NewInterfacesClient(subscriptionID)
	vmPublicIP:=network.NewPublicIPAddressesClient(subscriptionID)
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)
	vmClient.Authorizer = authorizer
	vmPublicIP.Authorizer = authorizer
	vmInterface.Authorizer = authorizer
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	results, err := vmClient.ListAllComplete(ctx)
	if err != nil {
		panic(err)
	}



	i:=0
	for results.NotDone() {
		vm := results.Value()
		i++
		fmt.Println(" ")
		fmt.Println("Vm no :",i )
		s1:= strings.Split(*vm.ID,"/")
		var subscription string
	    var resourceGroup string
	    var profile string
		var Name string 
		
		subscription = s1[2]
		resourceGroup = s1[4]
		profile = s1[6]
		Name = s1[8]
	    fmt.Println("Subscription ID : ",subscription)
	    fmt.Println("ResourceGroup : ",resourceGroup)
	    fmt.Println("profile  : ",profile)
	    fmt.Println("VMName : ",Name)
		for key, value := range vm.Tags {
			fmt.Println("Tags :", key, " :", *value)
		}
		fmt.Println("Location : ",*vm.Location)
		fmt.Println("Vm AdminUsername :", *vm.VirtualMachineProperties.OsProfile.AdminUsername)
		fmt.Println("Vm ID :", *vm.VirtualMachineProperties.VMID)
		fmt.Println("Size : ",vm.VirtualMachineProperties.HardwareProfile.VMSize)
		fmt.Println("Os : ",vm.VirtualMachineProperties.StorageProfile.OsDisk.OsType)
		vmnet := *vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces

		for i:=0;i<len(vmnet);i++{
			s2:=strings.Split(*vmnet[i].ID,"/")
			fmt.Println("vm network Interface ID :",s2[8])
			networkInterface := interfaceID(*vmnet[i].ID)
			interfaces,err:= vmInterface.Get(ctx,resourceGroup,networkInterface,"")
			if err != nil {
				panic(err)
			}
			interfaceinfo :=*interfaces.InterfacePropertiesFormat.IPConfigurations
			interfID := *interfaceinfo[0].InterfaceIPConfigurationPropertiesFormat
			fmt.Println("IP configuration :",*interfaceinfo[0].Name)
			if interfID.PrivateIPAddress!=nil {
			fmt.Println("PrivateIpaddress :",*interfID.PrivateIPAddress)
			}else{
				fmt.Println("PrivateIPaddress : -")
			}
			if interfID.PublicIPAddress!=nil&&interfID.PublicIPAddress.ID!=nil {
				ID:=strings.Split(*interfID.PublicIPAddress.ID,"/")
				fmt.Println("PublicIPAddress ID : ",ID[8])
				VmIP,err := vmPublicIP.Get(ctx,resourceGroup,ID[8],"")
				if err != nil {
					panic(err)
				}
				if VmIP.PublicIPAddressPropertiesFormat!=nil && VmIP.PublicIPAddressPropertiesFormat.IPAddress!=nil{
					fmt.Println("PublicIPAddress : ", *VmIP.PublicIPAddressPropertiesFormat.IPAddress)
				}else{
					
					fmt.Println("PublicIPAddress : -")
				}
				if VmIP.PublicIPAddressPropertiesFormat!=nil && VmIP.PublicIPAddressPropertiesFormat.DNSSettings!=nil{
					fmt.Println("DNS :", *VmIP.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn)
				}else{
					fmt.Println("DNS  : -")
				}
				
		   }else{
					fmt.Println("DNS : -")
					fmt.Println("PublicIPAddress ID: - ")
				
				}
		    if interfID.Subnet!=nil {
				ID := strings.Split(*interfID.Subnet.ID,"/")
				fmt.Println("Virtual Network/Subnet:",ID[8],"/",ID[10] )
				}else{
					fmt.Println("Virtual Network/Subnet:  -")
				}
		}
		
		if err := results.Next(); err != nil {
			panic(err)
		}
	}
        

}

