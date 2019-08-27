package chord

func remotePut(server *Info, value string) error {
	con, err := setupConnection(server)

	if err == nil {
		err = con.Call("Node.Put", value, nil)
	}

	return err
}

func remoteGet(server *Info, key string) (string, error) {
	var res string

	con, err := setupConnection(server)
	if err == nil {
		err = con.Call("Node.Get", key, &res)
	}

	return res, err
}

func remoteDelete(server *Info, key string) error {
	con, err := setupConnection(server)

	if err == nil {
		err = con.Call("Node.Delete", key, nil)
	}

	return err
}

func remoteSucc(server *Info) (*Info, error) {
	var res Info

	con, err := setupConnection(server)
	if err == nil {
		err = con.Call("Node.Successor", struct{}{}, &res)
	}

	return &res, err
}

func remotePred(server *Info) (*Info, error) {
	var res Info

	con, err := setupConnection(server)
	if err == nil {
		err = con.Call("Node.Predecessor", struct{}{}, &res)
	}

	return &res, err
}

func remoteNotify(server *Info, newPred Info) error {
	con, err := setupConnection(server)

	if err == nil {
		err = con.Call("Node.Notify", newPred, nil)
	}

	return err
}

func remoteFindPred(server *Info, id string) (*Info, error) {
	var pred Info

	con, err := setupConnection(server)
	if err == nil {
		err = con.Call("Node.FindPred", id, &pred)
	}

	return &pred, err
}

func remoteFindSucc(server *Info, id string) (*Info, error) {
	var succ Info

	con, err := setupConnection(server)
	if err == nil {
		err = con.Call("Node.FindSucc", id, &succ)
	}

	return &succ, err
}
