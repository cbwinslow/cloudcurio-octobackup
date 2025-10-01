import os
from crewai import Crew, Process
from agents import code_analyst, documentation_specialist, qa_engineer
from tasks import code_analysis_task, documentation_task, qa_review_task

# Create the documentation directory if it doesn't exist
docs_dir = "../docs"
if not os.path.exists(docs_dir):
    os.makedirs(docs_dir)

# Create the crew
code_review_crew = Crew(
    agents=[code_analyst, documentation_specialist, qa_engineer],
    tasks=[code_analysis_task, documentation_task, qa_review_task],
    process=Process.sequential,  # Execute tasks in sequence
    verbose=2,  # You can set it to 1 or 2 for different logging levels
)

# Run the crew
if __name__ == "__main__":
    print("Starting CloudCurio Code Review Crew...")
    result = code_review_crew.kickoff()
    print("Code Review Crew completed.")
    print("Result:", result)