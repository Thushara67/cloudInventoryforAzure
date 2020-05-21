package main

import (
	"github.com/Thushara67/cloudInventoryforAzure/azurelib"
	"fmt"
	"time"
	"context"
)

const subscriptionID = "282160c0-3c83-43f1-bff1-9356b1678ffb"

func main(){
	clients  := azurelib.GetNewClients(subscriptionID)
	clients  =  azurelib.AuthorizeClients(clients)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	
	s,err:=azurelib.GetallVMS(clients, ctx)
	if err != nil {
		panic(err)
	}
	
	for i:=0;i<len(s);i++{
		fmt.Println(" ")
		fmt.Println("vm no",i)
		subscription,err:=azurelib.GetVMSubscription(s[i])
		if err!=nil{
			fmt.Println("subscription ID : -" )
		}else{
			fmt.Println("subscription ID :",subscription)
		}
		resourceGroup,err := azurelib.GetVMResourcegroup(s[i])
		if err != nil {
			fmt.Println("ResourceGroup : -")
		}else{
		fmt.Println("Resourcegroup :",resourceGroup)
		}
		location,err := azurelib.GetVMLocation(s[i])
		if err != nil {
			fmt.Println("Location : -")
		}else{
		fmt.Println("Location :",location)
		}
		Name,err := azurelib.GetVMname(s[i])
		if err != nil {
		   fmt.Println("VM Name : -")	
		}else{
			fmt.Println("VM Name :",Name)
		}
		Tags,err := azurelib.GetVMTags(s[i])
		if err != nil {
			fmt.Println("Tags : -")
		}else{
			for key, value := range Tags {
				fmt.Println("Tags :", key, " :", *value)
			}
		}
		
		os:= azurelib.GetVMOsType(s[i])
		
		fmt.Println("Os :",os)
		Vmtype:= azurelib.GetVMSize(s[i])
		
		fmt.Println("Size :",Vmtype)

		networkinterface,err := azurelib.GetVmnetworkinterface(s[i])
		if err!=nil{
			fmt.Println("networkinterface : -")
			fmt.Println("IPConfiguration : -")
			fmt.Println("privateIPaddress : -")
			fmt.Println("publicIPname : -")
			fmt.Println("publicIPaddress : -")
			fmt.Println("virtualnetwork/subnet : -")
			fmt.Println("DNS : -")
			
		}else{ 
		   fmt.Println("networkinterface :",networkinterface)
		   privateipaddress,IPConfiguration,err:= azurelib.GetPrivateIP(clients, 
			ctx, resourceGroup, networkinterface, 
			"")
			if err!=nil{
				fmt.Println("privateIPaddress : -")
				fmt.Println("IPConfiguration :",IPConfiguration)
			}else{
				fmt.Println("privateIPaddress :",privateipaddress)
				fmt.Println("IPConfiguration :",IPConfiguration)
			}
		   publicIPname,err:=azurelib. GetPublicIPAddressID(clients,ctx, resourceGroup, networkinterface,"")
		   if err!=nil{
			 fmt.Println("publicIPname : -")
			 fmt.Println("publicIPaddress : -")
			 fmt.Println("DNS : -")
		   }else{
			 fmt.Println("publicIPname :",publicIPname)
			 publicipaddress,err := azurelib.GetPublicIPAddress(clients,ctx,resourceGroup, publicIPname ,"")
			 if err!=nil{
				fmt.Println("publicIPaddress: -")
			 }else{
				fmt.Println("publicIPaddress :",publicipaddress)
			 }
			 DNS,err := azurelib.GetDNS(clients,ctx,resourceGroup, publicIPname ,"")
			 if err!=nil{
				fmt.Println("DNS: -")
			 }else{
				fmt.Println("DNS :",DNS)
			 }
		   }
		   virtualnet,err:=azurelib.GetSubnetandvirtualnetwork(clients,ctx,resourceGroup,networkinterface,"")
		   if err!=nil{
			fmt.Println("virtualnetwork/subnet : -")
		   }else{
			fmt.Println("virtualnetwork/subnet : ",virtualnet)
		   }
		
	   }
	}
}
