# Robot32 - NB STILL IN PRE ALPHA - WORK IN PROGRESS - USE AT YOUR OWN RISK

The project has dual purpose:
* A - To contribute to open source AI ecosystem (by creating a frontent for LLMs)
* B - To create an AI that is good in electronics and robotics (mostly by RAG - retrieval augmented generation)

![Image of project UI](https://github.com/mlaak/robot32/blob/main/html/openscreen.png?raw=true)


## Archidecture

![Image of project archidecture](https://github.com/mlaak/robot32/blob/main/doc/r32diagram3.png?raw=true)


### Experts

Experts in this system function similarly to lightweight microservices. While they typically operate under the same Apache server rather than in separate containers, they maintain independence from both each other and the main system. Their primary requirement is access to configuration files.

This architecture aims to enhance the system's modularity. To reduce code duplication, commonly used logic is centralized in libraries, which can be found at https://github.com/mlaak/robot32lib. For instance, GPTlib handles the complexities associated with Large Language Models (LLMs).

Each expert maintains its own copies of the libraries it requires. These libraries are acquired either through Composer or through a custom system (which will be discussed in more detail later).

This approach allows for greater flexibility and easier maintenance, as experts can be developed, updated, or replaced independently without affecting the entire system. It also facilitates easier scaling and potential future containerization if needed.

An example of an expert is the **illustrator** (html/experts/illustrator), which generates illustrative images to accompany user queries, primarily for decorative purposes. By default, it attempts to connect to https://fal.ai to execute a picture generation model called 'SDXL-Turbo' on their servers. However, if the generation process fails due to fal.ai's occasional stability issues, the illustrator service selects the most fitting pre-generated image from its library based on the user query. This approach embodies microservice principles of independence and fault tolerance.

Experts have integration tests available, such as those located in the html/experts/illustrator/tests/integration directory. When you execute tools/test_all.php, it runs all the integration tests for the experts as well as the unit tests for the reverse proxy written in Golang.

### Reverse proxy (rate limiter)

Given that AI inference, including text and image generation, is computationally intensive, it is crucial to prevent system abuse. To achieve this, we utilize a limiting reverse proxy positioned in front of the experts who utilize the costly AI models. This proxy oversees user sessions and monitors usage, imposing both request and character count limits on a per-minute, hourly, and daily basis.

The proxy is written in Golang.


### UI



The current user interface (UI) is implemented in pure JavaScript, without the use of frameworks like React or Angular. This approach was partly chosen to optimize the main page's loading speed. However, In the near future, I plan to develop a React version of the UI to compare performance differences. For now, the UI remains in its pure JavaScript form and is somewhat messy.

To compile the HTML pages from their component parts located in `html/ui_parts`, please run the `tools/compile_html.php` script.

The UI offers two main features:
* 1. Interaction with the project's primary AI, which strives to excel in electronics and robotics.
* 2. Communication with various open-source and proprietary large language models (LLMs) such as Mixtral, ChatGPT, and Claude. This is facilitated through integration with OpenRouter.ai. Users need an OpenRouter account to access the more advanced (and costly) models.

Additionally, the UI enhances the user experience by displaying decorative images alongside user queries and AI responses. For users logged in via Google or GitHub, the UI also stores conversation history with the AI. These conversations are encrypted and stored locally in the browser, with the encryption key securely held on the server under the user's account. This ensures that the locally stored conversations become inaccessible upon user logout, maintaining data privacy.




















The UI is currently written in pure javascript (no react nor angular). This partly done to speed up the load time of the main page. I plan to soon make a version of UI in react and compare the speed difference. But for now it is pure javascript and a bit messy :)

Run tools/compile_html.php to compile the html pages from their parts (html/ui_parts). 

The UI lets the user do 2 main things:
 * 1) Talk to the main AI of this project. This AI strives to have enhanced capabilities in the areas of electronics and robotics
 * 2) Talk to a number of open and propiatry LLM models (like Mixtral or ChatGPT or Claude). This is possible thanks to integration with openrouter.ai. The user needs an openrouter account to access the more advanced (costly models).

Besides this, the UI displays nice decorative images alongside user queryes and AI anserws. If the user is logged in with Google or Github, then UI also stores the conversations that the user has with the AI. The data is stored in the browser in an encrypted form. The encryption key is stored in the server under the users account. This ensures that if the user loggs out then the locally stored conversations will become inaccessible.  











<!---


Based on the provided plan, it seems that the student's project, named Robot32, is primarily focused on developing a website (Robot32.com) that features a helpful AI, particularly in the field of technology, robotics, and automation. The AI will be built using open-source large language models (LLMs) from Mistral, such as Mistral 7b, Mixtral 8x7b, and Mixtral 8x22b.

The website will have a chat interface, allowing users to interact with the AI. The AI's behavior will be customized using Retrieval-Augmented Generation (RAG) and, in the future, possibly Lora training. These techniques help the AI to access and utilize relevant information during conversations, improving its ability to provide accurate and helpful responses.

While the project does not explicitly mention creating a physical robot for users to physically interact with, the AI on the Robot32 website will be able to provide guidance and resources for building hardware components, such as 'arms' and 'legs' for a robot. This way, the AI can assist users in creating their own physical robots by providing information and instructions.

Overall, the student's project aims to create a valuable and engaging web-based AI focused on technology, robotics, and automation, with a strong emphasis on open-source and customizable features.

Robot32 AI (Mixtral 8x7b)

-->