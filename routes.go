package tableapi

func (s *Server) routes() {

	// api methods
	s.router.HandleFunc("POST /records/{tablename}/", s.handleRecordID())       //  POST
	s.router.HandleFunc("PATCH /records/{tablename}/{id}/", s.handleRecordID()) // PATCH on /records/{id}
	s.router.HandleFunc("GET /records/{tablename}/", s.handleRecords())         // GET, POST

	// web interface methods
	s.router.HandleFunc("/readtable/{tablename}/", s.handleReadTable())
	s.router.HandleFunc("/editrow/{tablename}/{id}/{idx}/", s.handleEditRow())
	s.router.HandleFunc("/readrow/{tablename}/{id}/{idx}/", s.handleReadRow())
	s.router.HandleFunc("/saverow/{tablename}/{id}/{idx}/", s.handleSaveRow())

	s.router.HandleFunc("/static/", s.handleStatic())
	s.router.HandleFunc("/", s.handleFrontEnd())

}
