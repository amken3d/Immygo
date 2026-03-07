# Amken AI Usage Policy

**Amken LLC** | Effective March 2026

---

## My Philosophy

At Amken, AI is a force multiplier for an experienced engineer, not a replacement for one. I am a solo founder with 20+ years of systems engineering experience, a graduate background in electronics engineering, and hands-on leadership of large-scale engineering programs. AI accelerates my execution. It does not drive it.

I am transparent about how I use AI because I believe the community deserves to know, and because I think my approach is worth explaining on its merits.

---

## What I Use

I use **Claude by Anthropic** as my primary AI assistant for software and firmware development.

---

## Where AI Assists Me

### Software & Firmware Tooling

AI assists me with code generation, scaffolding, and implementation of software tools, desktop applications, and firmware for microcontrollers. This includes:

- Generating boilerplate and scaffolding from a well-defined architecture
- Implementing modules where the design and interfaces have already been specified
- Writing and refining documentation and README files
- Debugging and iterating on errors during development

### Documentation

AI assists with drafting READMEs, API documentation, and technical write-ups. All documentation is reviewed and edited by me to ensure it accurately reflects the actual implementation and my intent.

---

## How I Work: The Process

AI-generated code does not ship without going through my engineering process. Here is how I actually work:

1. **Architecture first.** Before AI touches anything, I design the system architecture. This includes module boundaries, data flow, interfaces, state management, and constraints. AI does not make architectural decisions — though I may use it to surface examples, suggestions, and help identify blind spots.

2. **Specification-driven prompting.** I prompt AI with precise, well-scoped specifications. Vague prompts produce vague code. Every prompt reflects a decision I have deliberately made.

3. **Review every line.** AI-generated code is read, understood, and evaluated by me before it is accepted. I do not merge code I do not understand.

4. **Test and verify.** Generated code is compiled, run, and tested against real hardware or environments. I do not assume correctness.

5. **Iterate with judgment.** When AI produces errors or suboptimal solutions, I diagnose the root cause myself and guide the correction. I do not blindly retry prompts.

6. **Domain knowledge governs.** My embedded systems, motion control, hardware domain expertise, and systems engineering experience is what catches what AI gets wrong. That expertise is not replaceable and is not outsourced.

---

## Hard Limits — What AI Never Does

These are absolute. No exceptions.

| Domain | Policy |
|--------|--------|
| **PCB & hardware design** | AI is never used. Schematics, layout, component selection, signal integrity, and power design are done entirely by me. |
| **Safety-critical firmware** | Motion control limits, fault detection, and any code path that can damage hardware or cause physical harm is human-written and reviewed by me. |
| **Hardware schematics** | Never AI-generated or AI-modified. |
| **Security & cryptographic code** | Never delegated to AI. Any security-sensitive implementation is written and audited by me. |

---

## What This Means for My Repositories

Projects in the Amken GitHub organization that involve AI-assisted development carry a notice in their README stating that AI was used and in what capacity.

The presence of AI assistance does not diminish the engineering behind a project. My hardware products — the Crea8 motion controller, TinyFOC ESC, and pick-and-place machine — are designed entirely by me without the use of AI. The software tools built around them may use AI assistance in implementation, but the systems they support do not.

---

## Why This Policy Exists

To sum it up in one word: transparency. I stand behind my products, my engineering, and my commitment to doing right by the communities I am part of. Publishing this policy is not a legal requirement — it is a reflection of who I am.

---

*Amken LLC — Springfield, VA USA*  
*github.com/amken3d*