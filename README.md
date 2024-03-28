# Magic-8-Ball

Welcome to the `magic-8-ball` repository! This command-line based application is a playful take on the classic Magic 8 Ball toy, brought to the digital era. Developed in Go, it's designed to help those new to the language and to LLM APIs get familiar with these technologies through a fun and interactive terminal user interface (TUI). Using the `tview` Go module for a responsive, platform-agnostic interface, users can ask the Magic 8 Ball any question and receive witty, humorous responses powered by LLM APIs.

Through this project I primarily learned Go and learned that I very much enjoy Go. This was also an excellent way for me to dig into the challenges that arise from working with LLMs programatically, as well as get practice with CI/CD while exploring new tools and platforms to build my knowledgebase.

## Getting Started

### Prerequisites

- Go installed on your machine.
- A Hugging Face API key, if you plan to run the project yourself.
- `gotty`, if you plan to run the project yourself as a webapp.

### Running Locally

1. Clone the repository to your local machine.
2. Install dependencies with `go get`
3. Navigate to the project directory and run `go run .` to start the application.

Note that when running the project yourself may run into issues with `glibc` due to the nature of drawing to a command-line interface (especially if you're hosting the program remotely) -- make sure that your system supports a wide range of C/C++ tools, platforms, and environments for best performance.

### Running on a Server with GOTTY

The project can also be run on a local or remote server, accessible via a web-based terminal emulator provided by `gotty`. This setup uses `xterm.js` under the hood for terminal emulation.

1. Build the Go project with `go build -o 8ball .`.
2. Run the binary with `gotty -w 8ball`.
3. Access the deployed binary using the localhost port `gotty` provides, and expose this port for remote access.

### Accessing the Deployed Application

The Magic 8 Ball application is deployed and available at: [https://8ball.aident.dev](https://8ball.aident.dev) for your enjoyment!

## Usage

After starting the application, simply type your question into the terminal and hit enter. Watch as the ASCII 8-ball performs a charming loading animation while it consults with a LLM to generate a response. The answer will then be displayed next to the 8-ball.

## Features

Currently, the Magic 8 Ball app offers a straightforward, delightful interaction with users. Future updates may include:

- Customizable response times for the Magic 8 Ball.
- A variety of animations to choose from for the loading sequence.
- Command-line arguments for enhanced usability and customization.

## LLM API Integration

This project utilizes the Hugging Face serverless inference API and Google's `gemma-7b-it` model for generating responses. If running the application locally, you will need to obtain an API key from Hugging Face.

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement."

Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Future Plans

- Integrating more LLM APIs for a diverse range of responses.
- Implementing user customization options for the TUI.
- Adding more interactive and engaging animations.
- Exploring additional CI/CD workflows to streamline development.

## Acknowledgments

- Special thanks to the developers of `tview`, `gotty`, and `xterm.js` for providing the tools that made this project possible.

Let's roll the magic dice and dive into the world of Go and LLMs together!
