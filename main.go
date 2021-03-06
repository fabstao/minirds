package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	rlista := ListarSvc()
	err := tpl.ExecuteTemplate(w, "dashboard.html", rlista)
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

//Listares is a struct...
type Listares struct {
	Resultado []string
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
	appName := nombre + "db"
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: nombre + "-depl",
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
							Name: "mariadbfabs",
							//Image: "gcr.io/fabs-cl-02/mariadbfabs",
							Image: "fabstao/mariadbfabs",
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
	fmt.Println("Creando K8s deployment para base de datos...")
	result, err := deploymentsClient.Create(deployment)
	if (err != nil) && (!strings.Contains(err.Error(), "already exists")) {
		return Creares{Resultado: "NULL (create deployment)", Error: err.Error()}
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

	servc := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	svc, err := servc.Get(appName, metav1.GetOptions{})
	switch {
	case err == nil:
		serviceSpec.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		//serviceSpec.Spec.LoadBalancerIP = svc.Spec.LoadBalancerIP
		_, err = servc.Update(serviceSpec)
		if (err != nil) && (!strings.Contains(err.Error(), "clusterIP")) {
			return Creares{Resultado: "NULL (servc->update)", Error: err.Error()}
		}
		fmt.Println("service updated: " + svc.GetName())
	case errors.IsNotFound(err):
		elsvc, errs := servc.Create(serviceSpec)
		if err != nil {
			return Creares{Resultado: "NULL (servc->Create)", Error: errs.Error()}
		}
		fmt.Println(time.Now().String()+" | service created: ", elsvc.GetClusterName())
	default:
		return Creares{Resultado: "NULL - default svc", Error: err.Error()}
	}
	//listasvc, err := servc.List(metav1.ListOptions{})
	svcres, err := servc.Get(appName, metav1.GetOptions{})
	if err != nil {
		fmt.Println(time.Now().String() + " | Error: " + err.Error())
	}
	elmapa := svcres.GetObjectMeta().GetSelfLink()

	//fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return Creares{Resultado: "Base de datos creada: " + result.GetObjectMeta().GetName() +
		" Svc:" + elmapa, Error: "OK"}

}

//ListarSvc is a function that lists DB services
func ListarSvc() Listares {
	// creates the in-cluster config
	os.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default.svc")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	config, err := rest.InClusterConfig()
	var salidaer []string
	if err != nil {
		salidaer[0] = "0"
		return Listares{Resultado: salidaer, Error: "Incluster config: " + err.Error()}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(time.Now().String() + " | Error: k8s config")
		salidaer[0] = "NULL"
		return Listares{Resultado: salidaer, Error: "Clientset config: " + err.Error()}
	}
	//lista, err := clientset.CoreV1().Services(apiv1.NamespaceDefault).List(metav1.ListOptions{})
	lista, err := clientset.CoreV1().Services(apiv1.NamespaceAll).List(metav1.ListOptions{})
	if err == nil {
		fmt.Println(time.Now().String()+" | Error creando lista: ", err.Error())
		salidaer[0] = ""
		return Listares{Resultado: salidaer, Error: "Error lista: " + err.Error()}
	}
	var slista []string
	for i, val := range lista.Items {
		slista[i] = val.String()
		fmt.Println(strconv.Itoa(i))
	}
	return Listares{Resultado: slista, Error: "OK"}
}

func int32p(i int32) *int32 {
	return &i
}
