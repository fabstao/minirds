package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
)

var host = "localhost"
var port = "3306"
var user = "gouser"
var password = "gopasswd"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("tps/*"))
}

//HandleError is a function...
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Registro is a function...
func Registro(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "registro.html", nil)
	HandleError(w, err)
}

//Login is a function...
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	HandleError(w, err)
}

//CreaServicio is a function...
func CreaServicio(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "crea.html", nil)
	HandleError(w, err)
}

//Index is a function...
func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	HandleError(w, err)
}

type ds1 struct {
	Test  string
	Value int
}

//Dashboard is a function...
func Dashboard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "dashboard.html", ds1{Test: "OK", Value: 1})
	HandleError(w, err)
}

//CreaDBI es una funcion ...
func CreaDBI(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	elres := crearDeplo(ps.ByName("nombre"))
	err := tpl.ExecuteTemplate(w, "dbi.html", elres)
	HandleError(w, err)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/registro.aspx", Registro)
	router.POST("/registro.aspx", Registro)
	router.GET("/crea.php", CreaServicio)
	router.POST("/crea.php", CreaServicio)
	router.GET("/login.aspx", Login)
	router.POST("/login.aspx", Login)
	router.GET("/dashboard.php", Dashboard)
	router.GET("/db/:nombre", CreaDBI)

	log.Fatal(http.ListenAndServe(":8800", router))
}

//Creares is a struct...
type Creares struct {
	Resultado string
	Error     string
}

//**************************************************+
func crearDeplo(nombre string) Creares {
	// creates the in-cluster config
	os.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default.svc")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	config, err := rest.InClusterConfig()
	if err != nil {
		return Creares{Resultado: "NULL", Error: "Incluster config: " + err.Error()}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return Creares{Resultado: "NULL", Error: "Clientset config: " + err.Error()}
	}

	//check deployments
	//var no2 int32 = int32(2)
	appName := nombre
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: nombre + "_depl",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32p(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "mariadbfabs",
							Image: "gcr.io/fabs-cl-02/mariadbfabs",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "mysql",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "changeme1st",
								},
							},
							/*VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      appName + "vol",
									MountPath: "/var/lib/mariadb",
								},
							}, */
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		return Creares{Resultado: "NULL", Error: err.Error()}
	}
	serviceSpec := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
		},
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeLoadBalancer,
			Selector: map[string]string{"app": appName},
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Protocol: apiv1.ProtocolTCP,
					Port:     3336,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(3306),
					},
				},
			},
		},
	}
	//pvclaim := clientset.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault)
	//pvc, err := pvclaim.Get(appName, metav1.GetOptions{})
	/* pvcspec := &apiv1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "vol",
			Labels: map[string]string{
				"app": appName,
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					"Storage": resource.Quantity{Format: resource.Format, d: intDecAmount(1000)},
				},
			},
		},
	} */

	servc := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	svc, err := servc.Get(appName, metav1.GetOptions{})
	switch {
	case err == nil:
		serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		serviceSpec.Spec.LoadBalancerIP = svc.Spec.LoadBalancerIP
		_, err = servc.Update(serviceSpec)
		if err != nil {
			return Creares{Resultado: "NULL", Error: err.Error()}
		}
		fmt.Println("service updated")
	case errors.IsNotFound(err):
		_, err = servc.Create(serviceSpec)
		if err != nil {
			return Creares{Resultado: "NULL", Error: err.Error()}
		}
		fmt.Println("service created")
	default:
		return Creares{Resultado: "NULL", Error: err.Error()}
	}
	//fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return Creares{Resultado: "Creado el deployment: " + result.GetObjectMeta().GetName(), Error: "OK"}

}

func int32p(i int32) *int32 {
	return &i
}
