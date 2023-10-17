import g4f

def main(*args, **kwargs):
    # Prompt
    prompt = """You are going to be given pseudo code. Your job is to transform it into actual code in the {{ LANGUAGE }} language.

    Text fields and comments should be in the same language as the pseudocode language. Try to give comments and print statements where possible.

    PSEUDO CODE:{{ CODE }}"""

    # Read file
    with open("work.pseudo", "r") as f:
        rawfile1 = f.read()
        rawfile2 = rawfile1
        
        lang = rawfile1.split("\n")[0].split(" ")[-1]
        code = rawfile2.split("\n")
        code.pop(0)
        code = "\n".join(code)
    
    # Make prompt
    prompt = prompt.replace("{{ LANGUAGE }}", lang)
    prompt = prompt.replace("{{ CODE }}", code)

    # Inference settings
    g4f.logging = False
    g4f.check_version = False
    bot = g4f.ChatCompletion

    # Inference
    response = bot.create(
        model="gpt-3.5-turbo",
        messages=[{"role": "user", "content": prompt}],
        stream=False,
    )


    print(response)
    
if __name__ == "__main__":
    main()