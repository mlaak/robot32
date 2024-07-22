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

An example of an expert is the *illustrator* (html/experts/illustrator), which generates illustrative images to accompany user queries, primarily for decorative purposes. By default, it attempts to connect to https://fal.ai to execute a picture generation model called 'SDXL-Turbo' on their servers. However, if the generation process fails due to fal.ai's occasional stability issues, the illustrator service selects the most fitting pre-generated image from its library based on the user query. This approach embodies microservice principles of independence and fault tolerance.






An example of an expert is the *illustrator* (html/experts/illustrator) which is a service that generates illustrative pictures to go with user queries (moslty for decorative purposes). By default it tries to connet to https://fal.ai to run a picture generation model 'SDXL-Turbo' on their servers. However if the generation fails (fal.ai has some stability problems) then the illustrator expert chooses a pre-generated image from it's library that it finds to be most appropriate for given user query. In that the expert encompasis the microservice ethos of being independent and fault tolerant.   







An expert is like a lightweight microservice. However these experts often run under same apache server (not in separate container, though they could). These 'lightwight microservices' are independent from others and independent from the main system - they only need access to conf files. Using this approach I am trying to make the system more modular. To minimize code duplication, a lot of re-occuring logic is put into libraries (under https://github.com/mlaak/robot32lib). For example GPTlib deals with the intricacies of dealing with LLMs (Large Language ai Models). Each expert has their own copies of the libraries they use (downloaded by composer or our own system - about that later).

An example of expert is the 'illustrator' (located html/experts/illustrator). It 




Based on the provided plan, it seems that the student's project, named Robot32, is primarily focused on developing a website (Robot32.com) that features a helpful AI, particularly in the field of technology, robotics, and automation. The AI will be built using open-source large language models (LLMs) from Mistral, such as Mistral 7b, Mixtral 8x7b, and Mixtral 8x22b.

The website will have a chat interface, allowing users to interact with the AI. The AI's behavior will be customized using Retrieval-Augmented Generation (RAG) and, in the future, possibly Lora training. These techniques help the AI to access and utilize relevant information during conversations, improving its ability to provide accurate and helpful responses.

While the project does not explicitly mention creating a physical robot for users to physically interact with, the AI on the Robot32 website will be able to provide guidance and resources for building hardware components, such as 'arms' and 'legs' for a robot. This way, the AI can assist users in creating their own physical robots by providing information and instructions.

Overall, the student's project aims to create a valuable and engaging web-based AI focused on technology, robotics, and automation, with a strong emphasis on open-source and customizable features.

Robot32 AI (Mixtral 8x7b)
