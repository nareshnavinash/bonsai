#!/usr/bin/env python3
"""
Bonsai Grammar Correction Benchmark

Tests Bonsai 1-bit models on 100 grammar/punctuation correction sentences
across multiple system prompt variants. Measures exact match rate, acceptable
corrections, verbose outputs, and latency.

Usage:
    # Test all models with default prompts
    python3 benchmarks/grammar-test.py

    # Test a specific model
    python3 benchmarks/grammar-test.py --model bonsai-4b

    # Test a specific prompt
    python3 benchmarks/grammar-test.py --prompt teacher

    # Test with a custom prompt
    python3 benchmarks/grammar-test.py --custom-prompt "Fix the grammar. Return only corrected text."

Prerequisites:
    - bonsai CLI installed and in PATH (or set BONSAI_BIN env var)
    - At least one model pulled (bonsai pull bonsai-4b)
"""

import json, time, urllib.request, re, os, subprocess, difflib, argparse, sys

BONSAI = os.environ.get("BONSAI_BIN", "bonsai")
PORT = int(os.environ.get("BONSAI_PORT", "8081"))
URL = f"http://127.0.0.1:{PORT}"
RESULTS_DIR = os.path.join(os.path.dirname(__file__), "results")

# ── System Prompts ──

PROMPTS = {
    "baseline": (
        "You are a grammar correction assistant. Fix grammar, spelling, and punctuation "
        "in the provided text. Return ONLY the corrected text, nothing else. Do not add "
        "explanations, notes, or markdown. Preserve the original meaning, tone, and formatting. "
        "If the text is already correct, return it as-is."
    ),
    "teacher": (
        "You are an English teacher helping a student. Read the following text and return "
        "a corrected version with proper grammar, spelling, and punctuation. Return only "
        "the corrected text, nothing else."
    ),
    "autocorrect": (
        "Fix grammar, spelling, and punctuation. Output the corrected text only. "
        "No explanations. No alternatives. No markdown. Minimal changes."
    ),
    "copy-editor": (
        "You are a copy editor. Fix only grammar, spelling, and punctuation errors in the "
        "text below. Make the smallest possible changes. Do not rewrite, rephrase, or improve "
        "style. Preserve the author's voice and word choices. Return only the corrected text."
    ),
    "fewshot": (
        "You are an English teacher. Fix grammar, spelling, and punctuation. "
        "Return only the corrected text.\n\n"
        "Examples:\n"
        "Input: Their going to the store.\n"
        "Output: They're going to the store.\n"
        "Input: I would of helped you.\n"
        "Output: I would have helped you.\n"
        "Input: The weather is more hotter.\n"
        "Output: The weather is hotter.\n"
        "Input: She dont like it.\n"
        "Input: Who's book is this.\n"
        "Output: Whose book is this?\n"
        "Input: I could care less about it.\n"
        "Output: I couldn't care less about it.\n"
        "Input: The principle gave a speech.\n"
        "Output: The principal gave a speech.\n"
        "Input: Irregardless of that we should go.\n"
        "Output: Regardless of that, we should go."
    ),
}

# ── 100 Test Cases ──
# Format: (input_with_errors, expected_correction)
# Covers 13 categories of common grammar/punctuation errors.

TEST_CASES = [
    # Subject-verb agreement (1-10)
    ("I goes to the store yesterday.", "I went to the store yesterday."),
    ("She dont like the movie at all.", "She doesn't like the movie at all."),
    ("The dogs is playing in the yard.", "The dogs are playing in the yard."),
    ("We was waiting for the bus.", "We were waiting for the bus."),
    ("She can sings very well.", "She can sing very well."),
    ("They doesn't want to come with us.", "They don't want to come with us."),
    ("The news are very surprising today.", "The news is very surprising today."),
    ("Each of the students have their own desk.", "Each of the students has their own desk."),
    ("Neither the teacher nor the students was prepared for it.", "Neither the teacher nor the students were prepared for it."),
    ("The committee have decided to postpone the meeting today.", "The committee has decided to postpone the meeting today."),

    # Pronoun errors (11-20)
    ("Him and me went to the park.", "He and I went to the park."),
    ("Me and him is best friends.", "He and I are best friends."),
    ("Me and her went shopping yesterday.", "She and I went shopping yesterday."),
    ("Him is the tallest person in class.", "He is the tallest person in class."),
    ("Between you and I this project is a disaster.", "Between you and me, this project is a disaster."),
    ("Who did you gave the report to this morning.", "To whom did you give the report this morning?"),
    ("I am interesting in learning new languages.", "I am interested in learning new languages."),
    ("She explained me the problem clearly yesterday.", "She explained the problem to me clearly yesterday."),
    ("He suggested me to take the job offer.", "He suggested that I take the job offer."),
    ("He made me to feel very uncomfortable yesterday.", "He made me feel very uncomfortable yesterday."),

    # Tense errors (21-30)
    ("She readed the book last night.", "She read the book last night."),
    ("He runned very fast in the race.", "He ran very fast in the race."),
    ("I seen him at the store last week.", "I saw him at the store last week."),
    ("She have been working here since 2020.", "She has been working here since 2020."),
    ("He gived me a present for my birthday.", "He gave me a present for my birthday."),
    ("I have went to Paris before.", "I have been to Paris before."),
    ("I did not went to the party last night.", "I did not go to the party last night."),
    ("She did not studied for the exam and failed badly.", "She did not study for the exam and failed badly."),
    ("Where did you went after the meeting ended yesterday.", "Where did you go after the meeting ended yesterday?"),
    ("She layed down on the couch and fell asleep quickly.", "She lay down on the couch and fell asleep quickly."),

    # Homophones and commonly confused words (31-40)
    ("Their going to the beach tommorow.", "They're going to the beach tomorrow."),
    ("Your the best person for this job position.", "You're the best person for this job position."),
    ("Its important too stay focused on you're goals.", "It's important to stay focused on your goals."),
    ("Who's book is this on the table.", "Whose book is this on the table?"),
    ("Their is a problem with the computer.", "There is a problem with the computer."),
    ("The affect of the new policy is very positive.", "The effect of the new policy is very positive."),
    ("Their going too the store for there supplies.", "They're going to the store for their supplies."),
    ("Whose coming to the party this weekend anyways.", "Who's coming to the party this weekend anyway?"),
    ("The principle of the school made a announcement.", "The principal of the school made an announcement."),
    ("The dog chased it's tail around the yard all day.", "The dog chased its tail around the yard all day."),

    # Would/could/should of (41-45)
    ("I would of helped if you asked me.", "I would have helped if you asked me."),
    ("I could of done better on the test today.", "I could have done better on the test today."),
    ("He should of went to the doctor last week.", "He should have gone to the doctor last week."),
    ("I could care less about there opinion on this matter.", "I couldn't care less about their opinion on this matter."),
    ("She could care less about what they think of her.", "She couldn't care less about what they think of her."),

    # Spelling errors (46-55)
    ("I need to practice my writting skills more often.", "I need to practice my writing skills more often."),
    ("He recieved the package yestarday afternoon finally.", "He received the package yesterday afternoon finally."),
    ("The goverment announced new enviromental regulations today.", "The government announced new environmental regulations today."),
    ("She definately wants to persue a career in medicene.", "She definitely wants to pursue a career in medicine."),
    ("The restarant has a very good enviroment and food.", "The restaurant has a very good environment and food."),
    ("alot of people showed up for the event today.", "A lot of people showed up for the event today."),
    ("I have no idear what your talking about right now.", "I have no idea what you're talking about right now."),
    ("Supposably the meeting was cancelled due to bad weather.", "Supposedly, the meeting was cancelled due to bad weather."),
    ("He pacifically asked for a refund on his purchase.", "He specifically asked for a refund on his purchase."),
    ("Irregardless of the weather we will go hiking tomorrow.", "Regardless of the weather, we will go hiking tomorrow."),

    # Capitalization and punctuation (56-65)
    ("I cant believe its already friday.", "I can't believe it's already Friday."),
    ("lets go to the movies tonight ok.", "Let's go to the movies tonight, OK."),
    ("i think we should leave now before its too late.", "I think we should leave now before it's too late."),
    ("the meeting is at 3pm on monday in room 5.", "The meeting is at 3 PM on Monday in room 5."),
    ("please send me the report asap thanks.", "Please send me the report ASAP. Thanks."),
    ("I like cooking my family and my pets.", "I like cooking, my family, and my pets."),
    ("We need to buy eggs milk bread and butter.", "We need to buy eggs, milk, bread, and butter."),
    ("However I think we should reconsider the plan.", "However, I think we should reconsider the plan."),
    ("Running through the park the dog chased the squirrel quickly.", "Running through the park, the dog chased the squirrel quickly."),
    ("After finishing dinner the movie started playing on TV.", "After finishing dinner, the movie started playing on TV."),

    # Comparative/superlative errors (66-70)
    ("The weather is more hotter than yesterday.", "The weather is hotter than yesterday."),
    ("She is more smarter than her brother.", "She is smarter than her brother."),
    ("He is very taller than his father.", "He is much taller than his father."),
    ("The movie was more better than I expected.", "The movie was better than I expected."),
    ("She did good on her exam today.", "She did well on her exam today."),

    # Countable/uncountable noun errors (71-75)
    ("The informations is very useful.", "The information is very useful."),
    ("We have less people than expected here.", "We have fewer people than expected here."),
    ("The amount of students are increasing every year.", "The number of students is increasing every year."),
    ("I have a informations for you about it.", "I have information for you about it."),
    ("The furnitures in the room is very old.", "The furniture in the room is very old."),

    # Preposition and article errors (76-80)
    ("I have been living here since five years.", "I have been living here for five years."),
    ("He has been here since three hours ago.", "He has been here for three hours."),
    ("He has less friends than his sister does.", "He has fewer friends than his sister does."),
    ("I am agree with your opinion completely.", "I agree with your opinion completely."),
    ("I am used to work late at night usually.", "I am used to working late at night usually."),

    # Gerund/infinitive errors (81-85)
    ("I look forward to hear from you soon.", "I look forward to hearing from you soon."),
    ("She avoid to speak in public situations.", "She avoids speaking in public situations."),
    ("She suggested me that I should apply for it.", "She suggested that I should apply for it."),
    ("I wish I can go to the concert tonight.", "I wish I could go to the concert tonight."),
    ("To get a good grade studying hard is necessary always.", "To get a good grade, studying hard is always necessary."),

    # Double negative and redundancy (86-90)
    ("He dont know nothing about it.", "He doesn't know anything about it."),
    ("She dont barely know anyone at the party.", "She barely knows anyone at the party."),
    ("The reason is because he was late to work today.", "The reason is that he was late to work today."),
    ("For all intensive purposes the project is complete now.", "For all intents and purposes, the project is complete now."),
    ("He did not only finish the work but also helped others.", "He not only finished the work but also helped others."),

    # Miscellaneous (91-100)
    ("The childrens toys was broken.", "The children's toys were broken."),
    ("The man who I saw him was very tall.", "The man who I saw was very tall."),
    ("The team are playing good today.", "The team is playing well today."),
    ("he said that \"he will come tomorrow\".", "He said that he will come tomorrow."),
    ("Everyone should bring their own lunch to the picnic.", "Everyone should bring their own lunch to the picnic."),
    ("If I was you I would take the opportunity now.", "If I were you, I would take the opportunity now."),
    ("She is one of those people who always helps others out.", "She is one of those people who always help others out."),
    ("There is alot of reasons why we should do this.", "There are a lot of reasons why we should do this."),
    ("The data shows that the trend is continuing upward.", "The data show that the trend is continuing upward."),
    ("I literally died laughing at his joke last night.", "I laughed so hard at his joke last night."),
]

CATEGORIES = [
    ("Subject-verb agreement", 0, 10),
    ("Pronoun errors", 10, 20),
    ("Tense errors", 20, 30),
    ("Homophones/confused words", 30, 40),
    ("Would/could/should of", 40, 45),
    ("Spelling errors", 45, 55),
    ("Capitalization/punctuation", 55, 65),
    ("Comparative/superlative", 65, 70),
    ("Countable/uncountable", 70, 75),
    ("Preposition/article", 75, 80),
    ("Gerund/infinitive", 80, 85),
    ("Double negative/redundancy", 85, 90),
    ("Miscellaneous", 90, 100),
]


# ── Helpers ──

def query(text, system_prompt, retries=3):
    """Send a correction request to the running bonsai server."""
    payload = json.dumps({
        "model": "bonsai",
        "messages": [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": text},
        ],
        "temperature": 0.1, "max_tokens": 256, "top_p": 0.9, "top_k": 20,
    }).encode()
    req = urllib.request.Request(
        f"{URL}/v1/chat/completions", data=payload,
        headers={"Content-Type": "application/json"},
    )
    for attempt in range(retries):
        start = time.time()
        try:
            resp = urllib.request.urlopen(req, timeout=60)
            data = json.loads(resp.read())
            ms = int((time.time() - start) * 1000)
            content = data["choices"][0]["message"]["content"].strip()
            content = re.sub(r"<think>.*?</think>", "", content, flags=re.DOTALL).strip()
            return content, ms
        except Exception as e:
            if attempt < retries - 1:
                time.sleep(2)
            else:
                return f"ERROR: {e}", int((time.time() - start) * 1000)


def start_model(model):
    """Stop any running server and start a new one for the given model."""
    subprocess.run([BONSAI, "stop"], capture_output=True)
    time.sleep(3)
    proc = subprocess.Popen(
        [BONSAI, "serve", model],
        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
    )
    for _ in range(90):
        try:
            r = urllib.request.urlopen(f"{URL}/health", timeout=3)
            if r.status == 200:
                break
        except Exception:
            pass
        time.sleep(1)
    for i in range(10):
        result, _ = query("Hello", "You are helpful.")
        if not result.startswith("ERROR"):
            print(f"  {model} ready (warmup {i + 1})")
            return proc
        time.sleep(2)
    print(f"  WARNING: {model} may not be ready")
    return proc


def score(orig, got, expected):
    """Score a model output against expected correction."""
    g, e = got.strip(), expected.strip()
    if g.startswith("ERROR:"):
        return "error"
    if len(g) > len(e) * 2.5 or g.count("\n") > 2:
        return "verbose"
    if g == e:
        return "exact"
    gn = re.sub(r"\s+", " ", g.lower()).strip()
    en = re.sub(r"\s+", " ", e.lower()).strip()
    if gn == en:
        return "near"
    sim = difflib.SequenceMatcher(None, g.lower(), e.lower()).ratio()
    if sim >= 0.90:
        return "acceptable"
    if g.lower().strip() == orig.lower().strip():
        return "unchanged"
    orig_sim = difflib.SequenceMatcher(None, g.lower(), orig.lower()).ratio()
    if sim > orig_sim and sim >= 0.70:
        return "partial"
    return "wrong"


def run_test(model, prompt_name, system_prompt):
    """Run 100 tests for a model+prompt combo."""
    results = []
    cats = {
        "exact": 0, "near": 0, "acceptable": 0, "partial": 0,
        "unchanged": 0, "wrong": 0, "verbose": 0, "error": 0,
    }
    total_ms = 0
    for i, (inp, exp) in enumerate(TEST_CASES):
        got, ms = query(inp, system_prompt)
        total_ms += ms
        s = score(inp, got, exp)
        cats[s] += 1
        results.append({
            "id": i + 1, "input": inp, "expected": exp,
            "got": got, "score": s, "ms": ms,
        })
        if (i + 1) % 25 == 0:
            en = cats["exact"] + cats["near"]
            print(
                f"    [{i + 1:3d}/100] exact+near={en} "
                f"accept={cats['acceptable']} wrong={cats['wrong']} "
                f"verbose={cats['verbose']} avg={total_ms // (i + 1)}ms"
            )
    n = len(TEST_CASES)
    en = cats["exact"] + cats["near"]
    os.makedirs(RESULTS_DIR, exist_ok=True)
    with open(os.path.join(RESULTS_DIR, f"{model}_{prompt_name}.json"), "w") as f:
        json.dump(results, f, indent=2)
    return {
        "model": model, "prompt": prompt_name,
        "exact": cats["exact"], "near": cats["near"], "exact_near": en,
        "acceptable": cats["acceptable"], "wrong": cats["wrong"],
        "verbose": cats["verbose"], "error": cats["error"],
        "unchanged": cats["unchanged"], "avg_ms": total_ms // n,
    }


def print_table(title, models, prompt_names, all_results, key):
    """Print a comparison table."""
    print(f"\n--- {title} ---\n")
    print(f"{'Model':<14}", end="")
    for p in prompt_names:
        print(f" {p:>14}", end="")
    print()
    print("-" * (14 + 15 * len(prompt_names)))
    for model in models:
        print(f"{model:<14}", end="")
        for p in prompt_names:
            r = next((x for x in all_results if x["model"] == model and x["prompt"] == p), None)
            val = r[key] if r else "N/A"
            suffix = "ms" if key == "avg_ms" else ""
            print(f" {val:>13}{suffix}", end="")
        print()


def main():
    parser = argparse.ArgumentParser(description="Bonsai Grammar Correction Benchmark")
    parser.add_argument("--model", type=str, help="Test a specific model (e.g. bonsai-4b)")
    parser.add_argument("--prompt", type=str, help="Test a specific prompt (e.g. teacher)")
    parser.add_argument("--custom-prompt", type=str, help="Test with a custom system prompt")
    args = parser.parse_args()

    models = [args.model] if args.model else ["bonsai-8b", "bonsai-4b", "bonsai-1.7b"]

    if args.custom_prompt:
        prompts_to_test = {"custom": args.custom_prompt}
    elif args.prompt:
        if args.prompt not in PROMPTS:
            print(f"Unknown prompt: {args.prompt}. Available: {', '.join(PROMPTS.keys())}")
            sys.exit(1)
        prompts_to_test = {args.prompt: PROMPTS[args.prompt]}
    else:
        prompts_to_test = PROMPTS

    prompt_names = list(prompts_to_test.keys())
    total_runs = len(models) * len(prompt_names)

    print(f"Bonsai Grammar Correction Benchmark")
    print(f"{len(TEST_CASES)} sentences x {len(prompt_names)} prompts x {len(models)} models = {total_runs} runs")
    print(f"Prompts: {', '.join(prompt_names)}")
    print(f"Models: {', '.join(models)}")

    all_results = []
    for model in models:
        print(f"\n{'=' * 60}")
        print(f"  MODEL: {model}")
        print(f"{'=' * 60}")
        proc = start_model(model)
        for pname in prompt_names:
            print(f"\n  Prompt: {pname}")
            r = run_test(model, pname, prompts_to_test[pname])
            all_results.append(r)
            print(
                f"  -> exact+near={r['exact_near']}  accept={r['acceptable']}  "
                f"wrong={r['wrong']}  verbose={r['verbose']}  avg={r['avg_ms']}ms"
            )
        proc.terminate()
        subprocess.run([BONSAI, "stop"], capture_output=True)
        time.sleep(3)

    # Print comparison tables
    print(f"\n\n{'=' * 80}")
    print("                         FINAL COMPARISON")
    print(f"{'=' * 80}")

    print_table("Exact + Near (strict correctness)", models, prompt_names, all_results, "exact_near")
    print_table("Verbose Outputs", models, prompt_names, all_results, "verbose")
    print_table("Avg Latency (ms)", models, prompt_names, all_results, "avg_ms")

    # Category breakdown for best combo
    best = max(all_results, key=lambda x: (x["exact_near"], -x["verbose"]))
    print(f"\n--- Best: {best['model']} + {best['prompt']} "
          f"(exact+near={best['exact_near']}, verbose={best['verbose']}) ---")

    # Save summary
    summary = {
        "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
        "prompts": {k: v for k, v in prompts_to_test.items()},
        "results": all_results,
    }
    os.makedirs(RESULTS_DIR, exist_ok=True)
    with open(os.path.join(RESULTS_DIR, "summary.json"), "w") as f:
        json.dump(summary, f, indent=2)

    print(f"\nDetailed results saved to: {RESULTS_DIR}/")


if __name__ == "__main__":
    main()
