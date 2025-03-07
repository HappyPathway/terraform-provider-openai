# Create a thread for analysis
resource "openai_thread" "analysis_session" {
  metadata = {
    session_type = "data_analysis"
    project      = "example"
  }

  tool_resources = {
    file_search = {
      vector_store_ids = [] # Empty list of strings
    }
  }
}
