package azurelib

import (
        "context"
        "testing"
        "time"
)


// TestGetAllSQLDBs tests the function GetallSQLDBs
func TestGetAllSQLDBs(t *testing.T) {
        if testing.Short() {
                t.Skip("Skipping test in short mode")
        }
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        subscriptions, err := GetAllSubscriptionIDs(ctx)
        if err != nil {
                t.Errorf("Unable to get subscriptionIDs: %v", err)
        }
        for key, subsID := range subscriptions {
                Dbs, err := GetAllSQLDBs(subsID)
                if err != nil {
                        t.Errorf("Failed to get databases for subscription: %s because %v", key, err)
                }
                t.Logf("Found %d databases in %s", len(Dbs), key)
        }
}
