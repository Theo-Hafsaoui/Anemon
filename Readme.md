![Anemon Project](https://github.com/user-attachments/assets/1399b964-5dfc-4dd5-b9ed-f333a3f768fe)  

# Anemon: CV Management and Generation  

Anemon is a simple tool designed to help you manage and generate tailored CVs.  
Crafting personalized CVs for every job application can be tedious and time-consuming, but Anemon aims to simplifies the process by focusing on two key features:  

1. **Language Support**  
   Anemon allows you to maintain CVs in multiple languages. Each language is represented by a separate directory (e.g., `eng` for English, `fr` for French),  
   where you can organize the corresponding sections of your CV.  

2. **Variant Support**  
   Variants make each CV unique by emphasizing specific keywords that align with a given application.  
   Anemon helps highlight these keywords and prioritizes relevant content to bring your skills to the forefront.  

If your CV becomes too lengthy, Anemon can automatically trim it by removing sections that lack or have insufficient emphasis on your chosen keywords.  
This allows you to be detailed in your Markdown source files without worrying about verbosity—irrelevant details will be excluded in the final CV.  

![Preview image](https://github.com/user-attachments/assets/19161e85-c666-40d1-978c-df4857075f13)  

### Getting Started

#### How to Create Your CV
To create your own CV, it is recommended to start with this repository, as the structure should be, I hope, intuitive.  
Two critical elements need your attention:  
1. **The CV Directory**  
   Each language has its own directory inside the `CV` folder. Within each language directory, you'll find Markdown documents representing  
   mandatory sections of the CV. These sections may vary in structure as some require more details (e.g., dates or years). The existing examples should suffice, but if you have any questions, don't hesitate to open an issue.  
2. **The `params.yaml` File**  
   This file contains your personal information and the profile(s) you want to generate CVs for. Each profile will create a new CV.

#### Using GitHub (Recommended)
To create your CV:  
1. **Fork** this repository.  
2. Update the `params.yaml` and CV file with your information.  
3. Once updated, the CI/CD pipeline will compile the CVs and publish them as artifacts on GitHub:  
   - Navigate to the **Actions** tab (next to the Pull Requests tab).  
   - Select the **Compile LaTeX Document** workflow.  
   - Inside the workflow, click on the job corresponding to your desired CV and download the `compiled-pdf.zip` file from the artifacts section.  

![Action Tab](https://github.com/user-attachments/assets/f15c7c71-022b-4bf2-b79d-2e5ef5f1e65e)  

#### Running Locally  
You can run Anemon locally using one of the following methods:  

1. **Using Docker**  
   Build the Docker image and run it using the provided `Makefile`. This setup creates a volume to simplify PDF extraction.  

   ```bash
   docker build -t anemon .  
   make run-docker  
   ```  

2. **Running Natively**  
   Alternatively, build and run Anemon natively using the `Makefile`:  
   ```bash
   ./anemon -g  
   ```  
   This approach requires a local installation of LaTeX and its necessary packages. Refer to the `Dockerfile` for the list of required dependencies and install their equivalents for your system.  

### Personalization  
You may want to customize the templates to fit your needs. Feel free to modify or create new templates, but ensure the `%something%` placeholders remain intact,  
as these serve as anchors for Anemon. If you want you can even share it by...

### Contribute  
Anemon has room for improvement and expansion. Contributions to enhance functionality or fix issues are welcome.  
Many parts of the project were initially designed as workarounds for personal use and can be improved.  
For example, the limited templates available reflect its original focus on personal needs.
