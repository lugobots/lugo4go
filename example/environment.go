package main

//
//func main() {
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//	})
//
//	go func() {
//		log.Fatal(http.ListenAndServe(":8080", nil))
//	}()
//
//	grpcServer := grpc.NewServer()
//
//	srv := lugo4go.GymEnv{}
//
//	proto.RegisterGymServer(grpcServer, srv)
//
//	lis, err := net.Listen("tcp", fmt.Sprintf(":2329"))
//	if err != nil {
//		log.Fatalf("failed on listen grpc port: %s", err)
//	}
//	log.Println("listening you")
//	err = grpcServer.Serve(lis)
//	log.Fatalf("stopped: %s", err)
//}
