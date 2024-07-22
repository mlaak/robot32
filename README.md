# Robot32 - NB STILL IN ALPHA - USE AT YOUR OWN RISK

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

Given that AI inference, including text and image generation, is computationally intensive, it is crucial to prevent system abuse. To achieve this, we utilize a limiting reverse proxy positioned in front of the experts who manage the costly AI models. This proxy oversees user sessions and monitors usage, imposing both request and character count limits on a per-minute, hourly, and daily basis.


Since AI inference (text and image generation) is computationally quite expensive it is important to prevent the abuse of the system. That is why there is a limiting reverse proxy that sits in front of the experts (experts that run computationally expensive AI models). The proxy handles user sessions and monitors usage. It has both request and capacity (nr of characters) limits for minute, hour and day.  












Based on the provided plan, it seems that the student's project, named Robot32, is primarily focused on developing a website (Robot32.com) that features a helpful AI, particularly in the field of technology, robotics, and automation. The AI will be built using open-source large language models (LLMs) from Mistral, such as Mistral 7b, Mixtral 8x7b, and Mixtral 8x22b.

The website will have a chat interface, allowing users to interact with the AI. The AI's behavior will be customized using Retrieval-Augmented Generation (RAG) and, in the future, possibly Lora training. These techniques help the AI to access and utilize relevant information during conversations, improving its ability to provide accurate and helpful responses.

While the project does not explicitly mention creating a physical robot for users to physically interact with, the AI on the Robot32 website will be able to provide guidance and resources for building hardware components, such as 'arms' and 'legs' for a robot. This way, the AI can assist users in creating their own physical robots by providing information and instructions.

Overall, the student's project aims to create a valuable and engaging web-based AI focused on technology, robotics, and automation, with a strong emphasis on open-source and customizable features.

Robot32 AI (Mixtral 8x7b)
