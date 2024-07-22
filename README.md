# Robot32 - NB STILL IN PRE ALPHA - WORK IN PROGRESS - USE AT YOUR OWN RISK

The project has dual purpose:
* A - To contribute to open source AI ecosystem (by creating a frontent for LLMs)
* B - To create an AI that is good in electronics and robotics (mostly by RAG - retrieval augmented generation)

![Image of project UI](https://github.com/mlaak/robot32/blob/main/html/openscreen.png?raw=true)


## 1 Archidecture

![Image of project archidecture](https://github.com/mlaak/robot32/blob/main/doc/r32diagram3.png?raw=true)


### 1.1 Experts

Experts in this system function similarly to lightweight microservices. While they typically operate under the same Apache server rather than in separate containers, they maintain independence from both each other and the main system. Their primary requirement is access to configuration files.

This architecture aims to enhance the system's modularity. To reduce code duplication, commonly used logic is centralized in libraries, which can be found at https://github.com/mlaak/robot32lib. For instance, GPTlib handles the complexities associated with Large Language Models (LLMs).

Each expert maintains its own copies of the libraries it requires. These libraries are acquired either through Composer or through a custom system (which will be discussed in more detail later).

This approach allows for greater flexibility and easier maintenance, as experts can be developed, updated, or replaced independently without affecting the entire system. It also facilitates easier scaling and potential future containerization if needed.

An example of an expert is the **illustrator** (html/experts/illustrator), which generates illustrative images to accompany user queries, primarily for decorative purposes. By default, it attempts to connect to https://fal.ai to execute a picture generation model called 'SDXL-Turbo' on their servers. However, if the generation process fails due to fal.ai's occasional stability issues, the illustrator service selects the most fitting pre-generated image from its library based on the user query. This approach embodies microservice principles of independence and fault tolerance.

Experts have integration tests available, such as those located in the html/experts/illustrator/tests/integration directory. When you execute tools/test_all.php, it runs all the integration tests for the experts as well as the unit tests for the reverse proxy written in Golang.

### 1.2 Reverse proxy (rate limiter)

Given that AI inference, including text and image generation, is computationally intensive, it is crucial to prevent system abuse. To achieve this, we utilize a limiting reverse proxy positioned in front of the experts who utilize the costly AI models. This proxy oversees user sessions and monitors usage, imposing both request and character count limits on a per-minute, hourly, and daily basis.

The proxy is written in Golang.


### 1.3 UI

The current user interface (UI) is implemented in pure JavaScript, without the use of frameworks like React or Angular. This approach was partly chosen to optimize the main page's loading speed. However, In the near future, I plan to develop a React version of the UI to compare performance differences. For now, the UI remains in its pure JavaScript form and is somewhat messy.

To compile the HTML pages from their component parts located in `html/ui_parts`, please run the `tools/compile_html.php` script.

The UI offers two main features:
* 1. Interaction with the project's primary AI, which strives to excel in electronics and robotics.
* 2. Communication with various open-source and proprietary large language models (LLMs) such as Mixtral, ChatGPT, and Claude. This is facilitated through integration with OpenRouter.ai. Users need an OpenRouter account to access the more advanced (and costly) models.

Additionally, the UI enhances the user experience by displaying decorative images alongside user queries and AI responses. For users logged in via Google or GitHub, the UI also stores conversation history with the AI. These conversations are encrypted and stored locally in the browser, with the encryption key securely held on the server under the user's account. This ensures that the locally stored conversations become inaccessible upon user logout, maintaining data privacy.

### 1.4 Notes about current implementation

Currently it uses openrouter.ai to do inference with LLM models and fal.ai to run image generation. Groq.com and stablehorde.net coming soon. Also fireworks.ai for custom LoRa trained LLM models is in the pipeline. 

In the future vast.ai gpu rental to train and self-host compact LLM models like Mistral Nemo.


## 2 Installation

Like mentioned, the project is not yet ready for primetime. However if you insist...

The project currently requires Linux or WSL under Windows, bash, curl, golang installation at `/usr/local/go/bin/go`, PHP>=7.4 with php libcurl installed.

Run `php tools/install_all.php` to install dependencies and create required directories. This also installs a version of composer under tools/composer. Besides composer, this project also uses a custom packege downloader. Composer does not support subpackages from monorepos, however some of my packages are simply not large enough to warrant their own repos and are thus serverd from a monorepo https://github.com/mlaak/robot32lib   

See `keys/README.md` for the api keys you need. This will get simpler soon for local installs but for now you need to go and get at-least openrouter key and hcaptcha key.

## 3 Testing

Run `php tools/test_all.php` This will run unit tests in goserver_reverse_proxy and integration tests for experts in `html/experts/*`

## 4 Running

To run dev server run `php run/run_dev.php` this will start a local php dev server at localhost:8000 and then starts the reverse proxy in front of it at localhost:8080

To run on production (warning not ready for it) start apache server to serve the html folder at localhost:8000 and run runserver_80.bash in goserver_reverse_proxy. Stay tuned for https (coming soon).

## 5 Current status

It is very much a work in progress. Can be a bit messy - but then again maybe not always good idea to too much polish things that might still change.

The robotics AI is not yet very smart. UI needs many improvements. Javascript code needs overhaul. Test covarage in all parts of the project needs to increase. Robot32lib currently has no tests yet, but some libraries there already would need them.

The install procedure needs to be tested and simplified. Should not need a login or captcha for local run. 

Ratelimiter configuration is currently hardcoded, needs to go into a conf file.

Https support for proxy should be ready very soon. 


<!---


Based on the provided plan, it seems that the student's project, named Robot32, is primarily focused on developing a website (Robot32.com) that features a helpful AI, particularly in the field of technology, robotics, and automation. The AI will be built using open-source large language models (LLMs) from Mistral, such as Mistral 7b, Mixtral 8x7b, and Mixtral 8x22b.

The website will have a chat interface, allowing users to interact with the AI. The AI's behavior will be customized using Retrieval-Augmented Generation (RAG) and, in the future, possibly Lora training. These techniques help the AI to access and utilize relevant information during conversations, improving its ability to provide accurate and helpful responses.

While the project does not explicitly mention creating a physical robot for users to physically interact with, the AI on the Robot32 website will be able to provide guidance and resources for building hardware components, such as 'arms' and 'legs' for a robot. This way, the AI can assist users in creating their own physical robots by providing information and instructions.

Overall, the student's project aims to create a valuable and engaging web-based AI focused on technology, robotics, and automation, with a strong emphasis on open-source and customizable features.

Robot32 AI (Mixtral 8x7b)

-->