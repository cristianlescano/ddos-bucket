package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func fetch(bucketInput string, ch chan bool, countCh chan int, sleep int) {
	// Configura la sesión de AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Cambia "tu-region" por la región de tu bucket S3
	})
	if err != nil {
		fmt.Println("Error creando sesión de AWS:", err)
		return
	}
	svc := s3.New(sess)
	key := "/archivo" + strconv.Itoa(rand.Intn(1000000000)) + ".txt" // Cambia "ruta/al/objeto/en/s3/archivo.txt" por la ruta de tu objeto en S3

	/*/ Archivo que deseas subir a S3
	file, err := os.Open("./archivo.txt") // Cambia "/ruta/de/tu/archivo/local/archivo.txt" por la ruta de tu archivo local
	if err != nil {
		fmt.Println("Error abriendo archivo:", err)
		return
	}
	defer file.Close()*/

	// Configura los parámetros para la operación PutObject
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketInput),
		Key:    aws.String(key),
		//Body:   file,
	}
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	// Incrementar el contador
	countCh <- 1
	// Realiza la operación PutObject
	_, err = svc.PutObject(params)
	if err != nil {
		tipoErr := strings.Split(strings.Split(err.Error(), "\n")[1], " ")[2]
		if tipoErr == "403," {
			ch <- true
		} else {
			//fmt.Println("Error subiendo archivo a S3:", err)
			ch <- false
		}
	}

}

func main() {
	// Definir la URL a la que se harán las solicitudes
	var bucket string
	fmt.Print("bucket: ")
	fmt.Scanf("%s\n", &bucket)

	// Definir la cantidad de goroutines (hilos) a abrir
	var numGoroutinesInput string
	fmt.Print("nro gorutinas: ")
	fmt.Scanf("%s\n", &numGoroutinesInput)

	numGoroutines, _ := strconv.Atoi(numGoroutinesInput)

	if numGoroutines < 1 {
		fmt.Println("El número de goroutines debe ser al menos 1")
		os.Exit(1)
	}

	var numSleepInput string
	fmt.Print("sleep (ms): ")
	fmt.Scanf("%s\n", &numSleepInput)

	numSleep, _ := strconv.Atoi(numSleepInput)

	// Crear un canal para comunicarse entre las goroutines y el hilo principal
	ch := make(chan bool)

	// Crear un canal para contar la cantidad de solicitudes exitosas
	countCh := make(chan int)

	// Iniciar 10 goroutines
	for i := 0; i < numGoroutines; i++ {
		go fetch(bucket, ch, countCh, numSleep)
	}

	// Mantener siempre 10 goroutines activas
	totalRequests := 0
	totalRequestsErr := 0
	for {
		select {
		case success := <-ch:
			if success {
				// Incrementar el contador de solicitudes exitosas
				totalRequests++
			} else {
				totalRequestsErr++
			}
			// Lanzar una nueva goroutine para reemplazarla
			go fetch(bucket, ch, countCh, numSleep)
		case count := <-countCh:
			// Imprimir el total de solicitudes exitosas
			fmt.Fprintf(os.Stdout, "\rTotal de solicitudes exitosas: %d, Errores: %d ", totalRequests+count, totalRequestsErr)
		}
	}
}

// update .syso
// $GOPATH/bin/rsrc -arch 386 -ico img/icon1.ico
// $GOPATH/bin/rsrc -arch amd64 -ico img/icon1.ico

// go build
