from typing import Dict, Any
import re

def validate_email_headers(
    dkim_signature: str = "",
    received_spf: str = "",
    from_header: str = "",
    return_path: str = "",
    reply_to: str = "",
    message_id: str = "N/A"
) -> Dict[str, Any]:
    ##Validates email headers with robust DKIM structure checks and risk scoring.

    def extract_domain(header: str) -> str:
        """Extracts domain from email addresses/headers."""
        match = re.search(r'@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})', header or "")
        return match.group(1).lower() if match else None

    def validate_dkim_structure(signature: str) -> Dict[str, Any]:
        ## Enhanced DKIM validation
        required_tags = {"v", "a", "d", "s", "b"}
        result = {
            "valid": False,
            "details": "Not checked",
            "missing_tags": []
        }

        try:
            if not signature:
                return result

            # Normalize signature and split tags
            normalized = re.sub(r"\s*;\s*", ";", signature.strip())
            parts = [p.strip().split("=", 1)[0].lower() 
                    for p in normalized.split(";") if "=" in p]
            
            missing = required_tags - set(parts)
            result["valid"] = not missing
            result["missing_tags"] = list(missing)
            result["details"] = ("Valid structure" if result["valid"] 
                                else f"Missing tags: {', '.join(missing)}")

        except Exception as e:
            result["details"] = f"Validation error: {str(e)}"

        return result

    results = {
        "message_id": message_id,
        "risk_score": 0,
        "risk_level": "low",
        "valid": False,
        "checks": {
            "dkim": {"valid": False, "details": "Not checked"},
            "spf": {"valid": False, "details": "Not checked"},
            "domain_alignment": {"valid": False, "details": "Not checked"},
            "reply_to_risk": {"valid": True, "details": "Not checked"}
        },
        "warnings": []
    }

    try:
        # --- DKIM Structure Validation ---
        if dkim_signature:
            dkim_check = validate_dkim_structure(dkim_signature)
            results["checks"]["dkim"] = {
                "valid": dkim_check["valid"],
                "details": dkim_check["details"]
            }
            if not dkim_check["valid"]:
                results["risk_score"] += 40
                results["warnings"].append("Invalid DKIM structure")

        # --- SPF Check ---
        if received_spf:
            results["checks"]["spf"]["valid"] = "pass" in received_spf.lower()
            results["checks"]["spf"]["details"] = "SPF passed" if results["checks"]["spf"]["valid"] else "SPF failed"
            if not results["checks"]["spf"]["valid"]:
                results["risk_score"] += 30

        # --- Domain Alignment ---
        from_domain = extract_domain(from_header)
        return_domain = extract_domain(return_path)
        
        if from_domain and return_domain:
            results["checks"]["domain_alignment"]["valid"] = (from_domain == return_domain)
            results["checks"]["domain_alignment"]["details"] = \
                f"Aligned: {from_domain}" if from_domain == return_domain \
                else f"From: {from_domain} ≠ Return-Path: {return_domain}"
            if not results["checks"]["domain_alignment"]["valid"]:
                results["risk_score"] += 20
        else:
            results["checks"]["domain_alignment"]["details"] = "Missing domains for comparison"
            results["risk_score"] += 10
            results["warnings"].append("Partial domain information")

        # --- Reply-To Risk ---
        reply_domain = extract_domain(reply_to)
        if reply_domain and from_domain:
            results["checks"]["reply_to_risk"]["valid"] = (reply_domain == from_domain)
            results["checks"]["reply_to_risk"]["details"] = \
                "Safe reply-to" if reply_domain == from_domain \
                else f"Risky reply-to: {reply_domain}"
            if not results["checks"]["reply_to_risk"]["valid"]:
                results["risk_score"] += 10

        # --- Final Scoring ---
        results["risk_score"] = min(results["risk_score"], 100)
        results["risk_level"] = (
            "high" if results["risk_score"] >= 70 else
            "medium" if results["risk_score"] >= 30 else
            "low"
        )
        results["valid"] = results["risk_score"] <= 20

    except Exception as e:
        results["warnings"].append(f"Validation error: {str(e)}")

    return results