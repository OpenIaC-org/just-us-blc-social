package chaincode

// branch test
import (
	"encoding/json"
	"fmt"

	//        "time"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	Owner          string   `json:"owner"`
	PostId         string   `json:"postId"`
	SharingHistory []string `json: "sharingHistory"`
	//PostDate string `json:"postDate"`
	//Poster   string `json:"poster"`
	//	Status		int	`json:"status"`
}

//type Post struct {
//	PostId string `json: "postId"`
//	SharingList string `json: "sharingList`
//}

type Profile struct {
	FollowedUsers    []string `json: followedUsers`
	Followers        []string `json:"followers"`
	PendingFollowers []string `json:"pendingFollowers"`
	Posts            []Asset  `json:"posts"`
	Username         string   `json: "username"`
}

//
//type memberType int
//
//const (
//	follower memberType = iota
//	blocked
//	pending
//)
//
////set sharing privileges for each individual user?
//type channelMember struct {
//	memberId string
//	membership memberType
//}
//

// init function neccessary for adding chaincode to peers. TODO: why neccessary, remove?
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	postList := make([]Asset, 0)
	postList = append(postList, Asset{
		SharingHistory: make([]string, 0),
		PostId:         "1",
		Owner:          "anders"})
	postList = append(postList, Asset{
		SharingHistory: make([]string, 0),
		PostId:         "2",
		Owner:          "anders",
	})
	userProfile := Profile{
		Followers:        make([]string, 0),
		FollowedUsers:    make([]string, 0),
		PendingFollowers: make([]string, 0),
		Posts:            postList,
		Username:         "anders",
	}
	userProfile2 := Profile{
		Followers:        make([]string, 0),
		FollowedUsers:    make([]string, 0),
		PendingFollowers: make([]string, 0),
		Posts:            make([]Asset, 0),
		Username:         "lars",
	}
	/* 	s.CreateProfile(ctx, "anders")
	   	s.CreateProfile(ctx, "abu")
	   	s.CreateProfile(ctx, "lars") */

	profileList := []Profile{userProfile, userProfile2}
	for _, asset := range profileList {
		profileJson, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(asset.Username, profileJson)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) GenerateSharingTree(ctx contractapi.TransactionContextInterface, userId string, postId string) ([]Asset, error) {
	iterableHistory, err := ctx.GetStub().GetHistoryForKey(userId)
	if err != nil {
		return make([]Asset, 0), err
	}
	var profileStates []Profile

	for iterableHistory.HasNext() {
		queryResult, err := iterableHistory.Next()
		if err != nil {
			return make([]Asset, 0), err
		}
		var profile Profile
		err = json.Unmarshal(queryResult.Value, &profile)
		if err != nil {
			return make([]Asset, 0), err
		}
		profileStates = append(profileStates, profile)
	}
	var postHistory []Asset
	for _, val := range profileStates {
		for _, value := range val.Posts {
			if value.PostId == postId {
				postHistory = append(postHistory, value)
			}
		}
	}
	return postHistory, nil
	//
	//	profileStatesJson, err := json.Marshal(profileStates)
	//	if err != nil {
	//		return make([]byte, 0), err
	//	}
	//
	//	return profileStatesJson, nil
}

// CreatePost issues a new asset to the world state with given details.
func (s *SmartContract) CreatePost(ctx contractapi.TransactionContextInterface, postId string, userId string, owner string) error {
	//id, err := ctx.GetClientIdentity().GetID()
	exists, err := s.PostExists(ctx, userId, postId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", postId)
	}

	sharingHistory := make([]string, 0)

	asset := Asset{
		Owner:          owner,
		PostId:         postId,
		SharingHistory: sharingHistory,
	}
	profile, err := s.ReadProfile(ctx, userId)
	if err != nil {
		return err
	}
	profile.Posts = append(profile.Posts, asset)
	profileJson, err := json.Marshal(profile)
	return ctx.GetStub().PutState(userId, profileJson)
}

//SharePost shares a post from another user to your own profile. Creator of post remains the original poster.
func (s *SmartContract) SharePost(ctx contractapi.TransactionContextInterface, owner string, postId string, userId string) error {
	//id, err := ctx.GetClientIdentity().GetID()
	exists, err := s.PostExists(ctx, userId, postId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the post in profile %s already exists", userId)
	}
	profile, _ := s.ReadProfile(ctx, owner)
	var post Asset
	var index int
	for i, item := range profile.Posts {
		if item.PostId == postId {
			post = item
			index = i
		}
	}
	post.SharingHistory = append(post.SharingHistory, userId)
	profile.Posts[index] = post
	profileJson, _ := json.Marshal(profile)
	err = ctx.GetStub().PutState(owner, profileJson)
	if err != nil {
		return err
	}
	// TODO: create sharing restriction after a certain number of sharing ops
	return s.CreatePost(ctx, postId, userId, owner)
}

// ReadPost returns the asset stored in the world state with given id.
func (s *SmartContract) ReadPost(ctx contractapi.TransactionContextInterface, userId string, postId string) (*Asset, error) {
	/* 	assetJSON, err := ctx.GetStub().GetState(id)
	   	if err != nil {
	   		return nil, fmt.Errorf("failed to read from world state: %v", err)
	   	}
	   	if assetJSON == nil {
	   		return nil, fmt.Errorf("the asset %s does not exist", id)
	   	}

	   	var asset Asset
	   	err = json.Unmarshal(assetJSON, &asset)
	   	if err != nil {
	   		return nil, err
	   	} */
	profileJson, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	var profile Profile
	json.Unmarshal(profileJson, &profile)
	for _, item := range profile.Posts {
		if item.PostId == postId {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("post not found")
}

// ReadProfile returns a profile complete with posts and other
func (s *SmartContract) ReadProfile(ctx contractapi.TransactionContextInterface, id string) (*Profile, error) {
	profileJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if profileJson == nil {
		return nil, fmt.Errorf("the profile %s does not exist", id)
	}

	var profile Profile
	err = json.Unmarshal(profileJson, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// PostExists returns true when asset with given ID exists in world state
func (s *SmartContract) PostExists(ctx contractapi.TransactionContextInterface, id string, postId string) (bool, error) {
	profileJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	var profile Profile
	json.Unmarshal(profileJson, &profile)
	if profileJson == nil {
		return false, fmt.Errorf("user %v does not exist", id)
	}
	for _, item := range profile.Posts {
		if item.PostId == postId {
			return true, nil
		}
	}
	return false, nil
}

//FollowProfile adds userId to followId's list of pending followers
func (s *SmartContract) FollowProfile(ctx contractapi.TransactionContextInterface, userId string, followId string) error {
	if userId == followId {
		return fmt.Errorf("users cannot follow themselves")
	}
	followProfile, err := s.ReadProfile(ctx, followId)
	if err != nil {
		return err
	}
	//check if user already sent request
	for _, val := range followProfile.PendingFollowers {
		if val == followId {
			return fmt.Errorf("Request already pending")
		}
	}
	for _, val := range followProfile.Followers {
		if val == followId {
			return fmt.Errorf("user %v already following user %v", userId, followId)
		}
	}
	userProfile, err := s.ReadProfile(ctx, userId)
	for _, val := range userProfile.FollowedUsers {
		if val == followId {
			return fmt.Errorf("user %v already following user %v", userId, followId)
		}
	}

	followProfile.PendingFollowers = append(followProfile.PendingFollowers, userId)
	profileJson, _ := json.Marshal(followProfile)
	return ctx.GetStub().PutState(followId, profileJson)
}

//AcceptFollower moves a user id from the list of pending followers to the list of followers
func (s *SmartContract) AcceptFollower(ctx contractapi.TransactionContextInterface, userId string, followerId string) error {
	profile, err := s.ReadProfile(ctx, userId)
	if err != nil {
		return err
	}
	for _, follower := range profile.Followers {
		if followerId == follower {
			return fmt.Errorf("this user already follows you!")
		}
	}
	for i, user := range profile.PendingFollowers {
		if user == followerId {
			profile.Followers = append(profile.Followers, user)
			profile.PendingFollowers[i] = profile.PendingFollowers[len(profile.PendingFollowers)-1]
			profile.PendingFollowers = profile.PendingFollowers[:len(profile.PendingFollowers)-1]
			profileJson, err := json.Marshal(profile)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(userId, profileJson)
			if err != nil {
				return err
			}

			followerProfile, err := s.ReadProfile(ctx, followerId)
			followerProfile.FollowedUsers = append(followerProfile.FollowedUsers, userId)
			profileJson, err = json.Marshal(followerProfile)
			if err != nil {
				return err
			}
			return ctx.GetStub().PutState(followerId, profileJson)
		}
	}

	return fmt.Errorf("User %v not found in pending follower list", followerId)
}

//create profile if profile does not already exist
func (s *SmartContract) CreateProfile(ctx contractapi.TransactionContextInterface, id string) error {
	//Abandoning using clientidentity for the moment, difficult to test
	//id, err := ctx.GetClientIdentity().GetID()
	exists, err := s.ProfileExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the profile %s already exists", id)
	}
	profile := Profile{
		Followers:        make([]string, 0),
		FollowedUsers:    make([]string, 0),
		PendingFollowers: make([]string, 0),
		Posts:            make([]Asset, 0),
		Username:         id,
	}
	profileJson, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, profileJson)
}

//func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, poster string) error {
//	exists, err := s.AssetExists(ctx, id)
//	if err != nil {
//		return err
//	}
//	if !exists {
//		return fmt.Errorf("the asset %s does not exist", id)
//	}

// overwriting original asset with new asset
//	asset := Asset{
//                Depth: 		depth,
//                Owner:          owner,
//                PostDate:       postDate,
//                Poster:         poster,
//                PostId:             id,
//                Status:         status,
//        }
//
//	assetJSON, err := json.Marshal(asset)
//	if err != nil {
//		return err
//	}
//
//	return ctx.GetStub().PutState(id, assetJSON)
//}
//
// DeleteAsset deletes an given asset from the world state.
/* func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.PostExists(ctx, id, postId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
} */

func (s *SmartContract) ProfileExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	profileJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return profileJson != nil, nil
}

func (s *SmartContract) GetAllPosts(ctx contractapi.TransactionContextInterface, userId string) ([]Asset, error) {
	profile, err := s.ReadProfile(ctx, userId)
	if err != nil {
		return nil, err
	}
	assetList := make([]Asset, 0)
	for _, user := range profile.FollowedUsers {
		followProfile, err := s.ReadProfile(ctx, user)
		if err != nil {
			return nil, err
		}
		assetList = append(assetList, followProfile.Posts...)
	}
	return assetList, nil
}

//figure out functions for changing state of assets such that they can be tracked
// TransferAsset updates the owner field of asset with given id in world state.
/* func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadPost(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
} */

// GetAllAssets returns all assets found in world state
/*  func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}  */
