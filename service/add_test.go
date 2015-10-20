package service

import (
	"testing"
)


//network:publishedPort:containerPort, network:containerPort, publishedPort:containerPort, containerPort
// in docker documentation, ports are named : hostPort:containerPort
// use those names in variables

func TestParsePublishedPortContainerPort(t *testing.T){

	portToTest := "80"

	// containerPort
	// expecting containerPorts["80/tcp"][0]{"80"}
	addPublish := []string{portToTest}

	containerPorts := make(map[string][]PortConfig)

	containerPorts = parsePublishedPort(addPublish)

	ContainerPortsLen := len(containerPorts)
	PortsConfigLen := len(containerPorts[portToTest + "/tcp"])
	
	if(ContainerPortsLen != 1){
		t.Error("Expecting a containerPort slice with 1 element, got ", ContainerPortsLen)
	}
	
	if(PortsConfigLen != 1){
		t.Error("Expecting a PortsConfig slice with 1 element, got ", PortsConfigLen)
	}

	if(containerPorts[portToTest+"/tcp"][0].PublishedPort != portToTest){
		t.Error("Expecting a PublishedPort on port " + portToTest + " , got ", containerPorts[portToTest + "/tcp"][0].PublishedPort)
	}
	
}

func TestParsePublishedPortPublishedPortContainerPort(t *testing.T){
	// publishedPort:containerPort
	// expecting containerPorts["80/tcp"][0]{"8080"}
	
	HostPort := "80"
	ContainerPort := "8080"
	
	addPublish := []string{HostPort + ":" + ContainerPort}

	containerPorts := make(map[string][]PortConfig)

	containerPorts = parsePublishedPort(addPublish)

	ContainerPortsLen := len(containerPorts)
	PortsConfigLen := len(containerPorts[HostPort + "/tcp"])

	if(ContainerPortsLen != 1){
		t.Error("Expecting a containerPort slice with 1 element, got ", ContainerPortsLen)
	}
	
	if(PortsConfigLen != 1){
		t.Error("Expecting a PortsConfig slice with 1 element, got ", PortsConfigLen)
	}

	if(containerPorts[HostPort+"/tcp"][0].PublishedPort != ContainerPort){
		t.Error("Expecting a PublishedPort on port " + ContainerPort + " , got ", containerPorts[HostPort + "/tcp"][0].PublishedPort)
	}
}

func TestParsePublishedPortNetworkContainerPort(t *testing.T){
	// network::containerPort
	// expecting containerPorts["8080/tcp"][0]{Network:"1.2.3.4" PublishedPort:"8080"}
	
	Network := "1.2.3.4"
	ContainerPort := "8080"

	addPublish := []string{Network + "::" + ContainerPort}

	containerPorts := make(map[string][]PortConfig)

	containerPorts = parsePublishedPort(addPublish)

	ContainerPortsLen := len(containerPorts)
	PortsConfigLen := len(containerPorts[Network+"/"+ContainerPort])

	if(ContainerPortsLen != 1){
		t.Error("Expecting a containerPort slice with 1 element, got ", ContainerPortsLen)
	}
	
	if(PortsConfigLen != 1){
		t.Error("Expecting a PortsConfig slice with 1 element, got ", PortsConfigLen)
		t.Error("Had object ", containerPorts)
	}

	if(containerPorts[ContainerPort+"/tcp"][0].PublishedPort != ContainerPort){
		t.Error("Expecting a PublishedPort on port " + ContainerPort + " , got ", containerPorts[ContainerPort+"/tcp"][0].PublishedPort)
	}

	if(containerPorts[ContainerPort+"/tcp"][0].Network != Network){
		t.Error("Expecting a Network " + Network + " , got ", containerPorts[ContainerPort+"/tcp"][0].Network)
	}
}

func TestParsePublishedPortNetworkPublishedPortContainerPort(t *testing.T){
	// network:publishedPort:containerPort
	// expecting containerPorts["1.2.3.4/80"][0]{Network:"1.2.3.4" PublishedPort:"8080"}
	Network := "1.2.3.4"
	ContainerPort := "8080"
	HostPort := "80"
	
	addPublish := []string{Network + ":" + HostPort +":" + ContainerPort}

	containerPorts := make(map[string][]PortConfig)

	containerPorts = parsePublishedPort(addPublish)

	ContainerPortsLen := len(containerPorts)
	PortsConfigLen := len(containerPorts[Network+"/"+ContainerPort])

	if(ContainerPortsLen != 1){
		t.Error("Expecting a containerPort slice with 1 element, got ", ContainerPortsLen)
	}
	
	if(PortsConfigLen != 1){
		t.Error("Expecting a PortsConfig slice with 1 element, got ", PortsConfigLen)
		t.Error("Had object ", containerPorts)
	}

	if(containerPorts[Network+"/"+ContainerPort][0].PublishedPort != HostPort){
		t.Error("Expecting a PublishedPort on port " + HostPort + " , got ", containerPorts[Network+"/"+ContainerPort][0].PublishedPort)
	}


	if(containerPorts[Network+"/"+ContainerPort][0].Network != Network){
		t.Error("Expecting a Network " + Network + " , got ", containerPorts[Network+"/"+ContainerPort][0].Network)
	}
}

