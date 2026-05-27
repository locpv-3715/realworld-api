package services

import (
	"realworld-api/testhelpers"
	"testing"
)

func TestGetProfile(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "target", "target@test.com", "pass")
	viewer := testhelpers.CreateTestUser(db, "viewer", "viewer@test.com", "pass")

	res, err := GetProfile(viewer.ID, "target")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Profile.Username != "target" {
		t.Errorf("expected username 'target', got '%s'", res.Profile.Username)
	}
	if res.Profile.Following {
		t.Error("expected following to be false")
	}
}

func TestGetProfile_WithFollowing(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	target := testhelpers.CreateTestUser(db, "target", "target@test.com", "pass")
	viewer := testhelpers.CreateTestUser(db, "viewer", "viewer@test.com", "pass")
	FollowUser(viewer.ID, target.Username)

	res, err := GetProfile(viewer.ID, "target")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.Profile.Following {
		t.Error("expected following to be true")
	}
}

func TestFollowUser_Service(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "target", "target@test.com", "pass")
	viewer := testhelpers.CreateTestUser(db, "viewer", "viewer@test.com", "pass")

	res, err := FollowUser(viewer.ID, "target")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.Profile.Following {
		t.Error("expected following to be true after follow")
	}
}

func TestFollowUser_Self(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "self", "self@test.com", "pass")
	_, err := FollowUser(user.ID, "self")
	if err == nil {
		t.Error("expected error when trying to follow yourself")
	}
}

func TestUnfollowUser_Service(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "target", "target@test.com", "pass")
	viewer := testhelpers.CreateTestUser(db, "viewer", "viewer@test.com", "pass")
	FollowUser(viewer.ID, "target")

	res, err := UnfollowUser(viewer.ID, "target")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Profile.Following {
		t.Error("expected following to be false after unfollow")
	}
}
