package testdata

import pb "github.com/haniz/byakugan/scanner/proto"

var SampleTasks = []*pb.ScanTask{
	{
		ID: "JWT-TEST-001",
		Target: &pb.Target{
			Url:    "http://api.example.com/users",
			Method: "GET",
		},
		AuthContext: &pb.AuthContext{
			Type: "jwt",
			Headers: map[string]string{
				"Authorization": "Bearer {{token}}",
			},
		},
		InsertionPoint: &pb.InsertionPoint{
			Location: "header.Authorization",
			Type:     "jwt_none",
		},
		Payload: "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJzdWIiOiJ0ZXN0In0.",
		Validation: &pb.Validation{
			SuccessConditions: &pb.SuccessConditions{
				StatusCodes:      []int32{200},
				ResponsePatterns: []string{"authenticated"},
			},
		},
	},
	// Add more sample tasks...
}
