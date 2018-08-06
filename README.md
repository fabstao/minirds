# minirds
# ***************************************************
# (C) 2018 Fabian Salamanca Dominguez
# ***************************************************

Basic mini RDS (Relational Database as a Service) based on MariaDB, using Jenkins for CI/CD, developed on Go

Based on official K8s.io client-go

Current version working only on Google's GKE.

TODO

Add list and delete functions
Show DB's access URL
Improve UI
Add user session capabilities to UI
Ask for support with clientset.CoreV1().Services(apiv1.NamespaceDefault).List(metav1.ListOptions{})
