# Office365-Validator
Office Validator is a multi-threaded Go-based tool for validating email addresses via API. It supports blacklist filtering, progress tracking, and detailed results, designed for efficiency and accuracy.

![image](https://github.com/user-attachments/assets/4c0cc071-bbab-48fe-9d73-62d4e7ca8914)



## Overview
Office Email Validator is a Go-based application designed to validate email addresses, specifically for Microsoft Office email services. It leverages multithreading for high performance and utilizes HTTP requests to verify email authenticity.

The project demonstrates robust error handling, efficient worker management, and sleek CLI presentation, making it an excellent tool for anyone working with email validation at scale.

---

## Features

- **Multi-threaded processing**: Efficiently processes up to 250 workers for large-scale email validations.
- **API Integration**: Utilizes Microsoft’s API endpoint for email validation.
- **Console Styling**: Clear and visually appealing CLI output with colors and formatting.
- **Blacklist Filtering**: Filters out blacklisted domains.
- **Robust Error Handling**: Handles network issues and retry logic for seamless processing.

---

## How It Works

1. The program reads an input file containing email addresses.
2. Blacklisted email domains are filtered out during preprocessing.
3. Each email is validated using a POST request to the Microsoft validation endpoint.
4. Results are displayed in the CLI and saved to an output file.

---

## Prerequisites

- **Go** (1.19 or later)
- Internet connection

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/pushkarup/office-email-validator.git
cd office-email-validator
```

2. Install dependencies:

```bash
go mod tidy
```

3. Build the application:

```bash
go build -o office-validator
```

---

## Usage

Run the application:

```bash
./office-validator
```

Follow the CLI instructions to provide the path to the email input file and the output file where validated results will be saved.

---

## Example Output

- **Valid Email:**
  ```
  [+]  |   VALID   | example@domain.com               | CRUX-CORE
  ```

- **Invalid Email:**
  ```
  [-]  |  INVALID  | example@invalid.com             | CRUX-CORE
  ```

---

## Skills Demonstrated

- **Golang Development**: Expertise in Go’s concurrency model and HTTP client.
- **API Integration**: Experience in consuming REST APIs with JSON payloads.
- **Multithreading**: Proficient in using Goroutines and channels for parallel processing.
- **Error Handling**: Implementing retry logic and fail-safe mechanisms.
- **CLI Design**: Creating intuitive and user-friendly CLI applications.

---

## Future Enhancements

- **Dockerization**: Create a Docker image for easy deployment.
- **Configurable Blacklist**: Add support for custom blacklists.
- **Support for Additional APIs**: Expand to validate emails across different services.

---

## Contact

**Author**: Pushkar Upadhyay 

**GitHub**: [GitHub Profile](https://github.com/pushkarup)  
**LinkedIn**: [LinkedIn Profile](https://linkedin.com/in/pushkar-upadhyay-381634315)

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.

---

## Contribution

Contributions are welcome! Feel free to open issues or submit pull requests to enhance the project.

---

Thank you for checking out the Office Email Validator. Don’t forget to star the repository if you found it helpful!


