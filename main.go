package main

import (
	"os/exec"
	"fmt"
	"log"
	"net/http"
)

var (
	script_path = "./create_sa_user.sh"
	kubecfg_path = "/tmp"
	default_namespace = "test"
)

func createNamespace(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	namespace := r.Form["namespace"][0]

	kubecfg_file := fmt.Sprintf("%s/%s", kubecfg_path, "k8s-config-" + namespace)
	cmd := []string{"create", namespace, kubecfg_path}
	out, err := exec.Command(script_path, cmd...).Output()
	if err != nil {
        fmt.Printf("Failed to execute command: %v, output: %s, Error: %v", cmd, out, err)
		fmt.Fprintf(w, fmt.Sprintf("Failed to create namespace: %s", namespace))
		return
    }
	fmt.Printf("Output of command: %v, output: %s", cmd, out)

	w.Header().Set("Content-Disposition", "attachment; filename=config")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	http.ServeFile(w, r, kubecfg_file)
}

func deleteNamespace(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, fmt.Sprint("I am sorry, but I don't know how to delete namespaces yet. Would you teach me that and submit PR at https://github.com/prgcont/workshop-namespaces please?"))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/create", createNamespace) // set router
	http.HandleFunc("/delete", deleteNamespace) // set router
	log.Fatal(http.ListenAndServe(":9090", nil))
}
