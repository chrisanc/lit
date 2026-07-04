# **lit: code linter**
### **What's 'lit'? ❓**
CLI Tool developed 100% in Golang in order to help you get quick information about your scripts in determined languages. Thus, it can be really helpful for big work enviroments in order to improve productivity or detect code anomalies, getting warnings about possible dangerous code.

I used hexagonal architecture on the project, learning about mantainable software architectures.

### **Top features 🔝**
With the available commands, you can scan your scripts and get useful feedback that helps you to improve your code.
The scanning does a search in your working directory looking for the scripts written on the supported languages to scan and calculates the cyclical complexity so you can get feedback based in the metrics configured in the json file.
It also fix your naming conventions with the flag --fix (this can make some mistakes, so use carefully).

### **Future features (plans)**
Currently (release v1.0.0), the 'scan' command without flags only scans methods/functions but I'll implement the class scanning so it also detects the bad naming conventions on the properties.

### **Supported languages for the code scan 💻**
Current supported languages: **Java, Python,  JavaScript, Go**

Incoming support for the next languages: **C, C++, C#, TypeScript**

### **Requirements 🧰**
- 64x C compiler installed (for the file scanner)
  
  That's the only requirement to run the project.

### **How can I set it up?**:
To set the project up in your machine and start scanning your projects, you must:
- Download the latest release .zip (which contains the executable and the configuration file).
- Unzip the file and set the directory into the program files of yours system.
- Create a enviroment variable pointing to the directory path of your system (so you can just write 'lit' to use it).
- Confirm the changes and start using lit.

### **Available commands 🌝**
- *lit files*: The brain and main command of this project. This command itself scans your whole repository and finds the scripts of the supported languages, scanning and looking for possible dangerous functions defined
  and variables that doesn't match with the naming convention defined.
  *Command flags*:
  - *--loc*: This flag allows you to know how much lines of code were written on which language and their percentage of the total lines.
  - *--fix*: This is a powerful flag which can fix most of the detected naming conventions into the one you wish in your project. This flag can make some mistakes and might not fix every name for safety reasons (code safety).

- *lit config*: With this, you can modify the **config.json** file with a friendly and easy interface. Currently, the only configuration is the current regex for the naming convention you're using.
  I would recommend to run this command before anything. (the default convention is **CamelCase/camelCase**). Also, you can personalize the minimum parameters, method size or complexity the scanner finds to 

  The available conventions are: **LowerCamelCase, UpperCamelCase, CamelCase, snake_case**.

### Why was this developed?
I developed ***Lit*** because I wanted to reinforce my Go knowledge by creating a useful and meaningful project that developers like me could use in their development projects to keep the code cleaner. Writing clean but efficient code is very important
because the code will always be read by developers and they must understand it. Developing this project I learned about go routines and how powerful they are, also, I was able to reinforce my Go knowledge and now I feel ready to start bigger projects.
