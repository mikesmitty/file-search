package gemini

import (
	"context"
	"testing"
)

func TestIntegration_ListStores(t *testing.T) {
	client, cleanup := newVCRClient(t, "list_stores")
	defer cleanup()

	ctx := context.Background()
	stores, err := client.ListStores(ctx)
	if err != nil {
		t.Fatalf("ListStores failed: %v", err)
	}

	t.Logf("Found %d stores", len(stores))
	for _, s := range stores {
		if s.Name == "" {
			t.Error("Store name is empty")
		}
		if s.DisplayName == "" {
			t.Error("Store display name is empty")
		}
	}
}

func TestIntegration_ResolveNames(t *testing.T) {
	client, cleanup := newVCRClient(t, "resolve_names")
	defer cleanup()

	ctx := context.Background()
	
	stores, err := client.ListStores(ctx)
	if err != nil {
		t.Fatalf("ListStores failed: %v", err)
	}
	if len(stores) == 0 {
		t.Skip("No stores available to test resolution")
	}

	targetStore := stores[0]
	
	id, err := client.ResolveStoreName(ctx, targetStore.DisplayName)
	if err != nil {
		t.Fatalf("Failed to resolve store name %q: %v", targetStore.DisplayName, err)
	}

	if id != targetStore.Name {
		t.Errorf("Resolved ID mismatch. Got %s, want %s", id, targetStore.Name)
	}
}

func TestIntegration_Query(t *testing.T) {
	client, cleanup := newVCRClient(t, "query_response")
	defer cleanup()

	ctx := context.Background()
	
	stores, err := client.ListStores(ctx)
	if err != nil {
		t.Fatalf("ListStores failed: %v", err)
	}
	if len(stores) == 0 {
		t.Skip("No stores available to test query")
	}

	resp, err := client.Query(ctx, "What is in this document?", stores[0].Name, "gemini-2.5-flash", "")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if resp == nil || len(resp.Candidates) == 0 {
		t.Fatal("Query returned no candidates")
	}
	
	t.Logf("Query response: %v", resp.Candidates[0].Content.Parts[0].Text)
}

func TestIntegration_ListFiles(t *testing.T) {
	client, cleanup := newVCRClient(t, "list_files")
	defer cleanup()

	ctx := context.Background()
	files, err := client.ListFiles(ctx)
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	t.Logf("Found %d files", len(files))
	for _, f := range files {
		if f.Name == "" {
			t.Error("File name is empty")
		}
		if f.DisplayName == "" {
			t.Error("File display name is empty")
		}
	}
}

func TestIntegration_GetStore(t *testing.T) {
	client, cleanup := newVCRClient(t, "get_store")
	defer cleanup()

	ctx := context.Background()
	stores, err := client.ListStores(ctx)
	if err != nil {
		t.Fatalf("ListStores failed: %v", err)
	}
	if len(stores) == 0 {
		t.Skip("No stores available to test GetStore")
	}

	targetStore := stores[0]
	store, err := client.GetStore(ctx, targetStore.Name)
	if err != nil {
		t.Fatalf("GetStore failed: %v", err)
	}

	if store.Name != targetStore.Name {
		t.Errorf("Store ID mismatch. Got %s, want %s", store.Name, targetStore.Name)
	}
	if store.DisplayName != targetStore.DisplayName {
		t.Errorf("Store DisplayName mismatch. Got %s, want %s", store.DisplayName, targetStore.DisplayName)
	}
}

func TestIntegration_GetFile(t *testing.T) {
	client, cleanup := newVCRClient(t, "get_file")
	defer cleanup()

	ctx := context.Background()
	files, err := client.ListFiles(ctx)
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}
	if len(files) == 0 {
		t.Skip("No files available to test GetFile")
	}

	targetFile := files[0]
	file, err := client.GetFile(ctx, targetFile.Name)
	if err != nil {
		t.Fatalf("GetFile failed: %v", err)
	}

	if file.Name != targetFile.Name {
		t.Errorf("File ID mismatch. Got %s, want %s", file.Name, targetFile.Name)
	}
}

func TestIntegration_GetDocument(t *testing.T) {
	client, cleanup := newVCRClient(t, "get_document")
	defer cleanup()

	ctx := context.Background()
	
	// 1. Find a store
	stores, err := client.ListStores(ctx)
	if err != nil {
		t.Fatalf("ListStores failed: %v", err)
	}
	if len(stores) == 0 {
		t.Skip("No stores available to test GetDocument")
	}
	store := stores[0]

	// 2. List documents in that store
	docs, err := client.ListDocuments(ctx, store.Name)
	if err != nil {
		t.Fatalf("ListDocuments failed: %v", err)
	}
	if len(docs) == 0 {
		t.Skip("No documents available in store to test GetDocument")
	}

	// 3. Get the first document
	targetDoc := docs[0]
	doc, err := client.GetDocument(ctx, targetDoc.Name)
	if err != nil {
		t.Fatalf("GetDocument failed: %v", err)
	}

	if doc.Name != targetDoc.Name {
		t.Errorf("Document ID mismatch. Got %s, want %s", doc.Name, targetDoc.Name)
	}
}