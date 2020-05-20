package main

import (
	"Azure/lib/azurelib"
	"fmt"
)

const subscriptionID = "282160c0-3c83-43f1-bff1-9356b1678ffb"

func main(){
	s,err:=azurelib.GetallVMS(subscriptionID)
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
			fmt.Println("publicIPname : -")
			fmt.Println("publicIPaddress : -")
			fmt.Println("virtualnetwork/subnet : -")
		}else{ 
		   fmt.Println("networkinterface :",networkinterface)
		   publicipname,err:=azurelib.GetPublicIPname(subscriptionID,resourceGroup,
			networkinterface)
		   if err!=nil{
			 fmt.Println("publicIPname : -")
			 fmt.Println("publicIPaddress : -")
		   }else{
			 fmt.Println("publicIPname :",publicipname)
			 publicipaddress,err := azurelib.GetPublicIpaddress(subscriptionID,resourceGroup,publicipname)
			 if err!=nil{
				fmt.Println("publicIPaddress: -")
			 }else{
				fmt.Println("publicIPaddress :",publicipaddress)
			 }
		   }
		   virtualnet,err:=azurelib.GetSubnetandvirtualnetwork(subscriptionID,resourceGroup,networkinterface)
		   if err!=nil{
			fmt.Println("virtualnetwork/subnet : -")
		   }else{
			fmt.Println("virtualnetwork/subnet : ",virtualnet)
		   }
		
	   }
	}
}