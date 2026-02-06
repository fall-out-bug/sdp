"""
Question definitions for @idea interview rounds.

Provides structured questions for Round 1 (critical) and Round 2 (deep dive).
"""

from typing import Any, Dict, List


class CriticalQuestions:
    """Critical questions for Round 1 of @idea interview."""

    @staticmethod
    def get_round_1_questions() -> List[Dict[str, Any]]:
        """Get critical questions for Round 1.

        Returns:
            List of question definitions compatible with AskUserQuestion
        """
        return [
            {
                "question": "What is the primary problem this feature solves?",
                "header": "Problem",
                "options": [
                    {
                        "label": "User pain point",
                        "description": "Addresses frustration or inefficiency",
                    },
                    {
                        "label": "Business requirement",
                        "description": "Enables new revenue or reduces cost",
                    },
                    {
                        "label": "Technical debt",
                        "description": "Improves maintainability or performance",
                    },
                    {
                        "label": "Competitive parity",
                        "description": "Matches competitor capabilities",
                    },
                ],
                "multiSelect": False,
            },
            {
                "question": "Who are the primary users of this feature?",
                "header": "Users",
                "options": [
                    {"label": "End users", "description": "Direct product users"},
                    {
                        "label": "Administrators",
                        "description": "System managers and ops teams",
                    },
                    {
                        "label": "Developers",
                        "description": "Engineering team integration",
                    },
                    {"label": "API consumers", "description": "External integrations"},
                ],
                "multiSelect": True,
            },
            {
                "question": "What is the technical approach?",
                "header": "Technical Approach",
                "options": [
                    {
                        "label": "Database-driven",
                        "description": "PostgreSQL, MySQL, MongoDB",
                    },
                    {"label": "In-memory", "description": "Redis, Memcached"},
                    {"label": "File-based", "description": "JSON, YAML, CSV"},
                    {"label": "External API", "description": "Third-party service"},
                ],
                "multiSelect": False,
            },
            {
                "question": "What is the risk level?",
                "header": "Risk",
                "options": [
                    {
                        "label": "Low",
                        "description": "Well-understood problem, proven solution",
                    },
                    {"label": "Medium", "description": "Some unknowns, but manageable"},
                    {"label": "High", "description": "New territory, high uncertainty"},
                ],
                "multiSelect": False,
            },
        ]

    @staticmethod
    def get_round_2_questions(ambiguous_fields: List[str]) -> List[Dict[str, Any]]:
        """Get deep dive questions based on ambiguous fields.

        Args:
            ambiguous_fields: List of fields that were ambiguous in Round 1

        Returns:
            List of follow-up questions
        """
        questions = []

        # Map ambiguous fields to deep dive questions
        field_questions = {
            "problem": [
                {
                    "question": "Can you describe a specific use case or scenario?",
                    "header": "Problem Clarification",
                    "options": [
                        {
                            "label": "Specific scenario",
                            "description": "I have a concrete example",
                        },
                        {
                            "label": "General pain point",
                            "description": "Broad problem area",
                        },
                    ],
                    "multiSelect": False,
                }
            ],
            "users": [
                {
                    "question": "Who specifically will use this feature?",
                    "header": "User Clarification",
                    "options": [
                        {"label": "Internal team", "description": "Our developers"},
                        {"label": "Customers", "description": "End customers"},
                        {"label": "Partners", "description": "External integrators"},
                    ],
                    "multiSelect": True,
                }
            ],
            "tech_approach": [
                {
                    "question": "What are the specific technical requirements?",
                    "header": "Technical Details",
                    "options": [
                        {
                            "label": "ACID transactions",
                            "description": "Strong consistency needed",
                        },
                        {
                            "label": "Eventually consistent",
                            "description": "High throughput, relaxed consistency",
                        },
                        {"label": "No persistence", "description": "Ephemeral data only"},
                    ],
                    "multiSelect": False,
                }
            ],
        }

        # Add questions for ambiguous fields
        for field_name in ambiguous_fields:
            if field_name in field_questions:
                questions.extend(field_questions[field_name])

        return questions
