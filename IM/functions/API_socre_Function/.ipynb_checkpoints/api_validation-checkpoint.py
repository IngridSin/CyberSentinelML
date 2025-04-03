from typing import Dict, Any
import re

def assess_email_risk(
    dkim_signature: str = "",
    received_spf: str = "",
    from_header: str = "",
    return_path: str = "",
    reply_to: str = "",
    message_id: str = "N/A"
) -> Dict[str, Any]:
    
    def extract_domain(header: str) -> str:
        match = re.search(r'@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})', header or "")
        return match.group(1).lower() if match else ""

    results = {
        "message_id": message_id,
        "valid": True,  # Innocent until proven guilty
        "risk_score": 0,
        "risk_level": "Low",
        "detailed_checks": {
            "dkim": {"status": "not_checked", "details": ""},
            "spf": {"status": "not_checked", "details": ""},
            "domain_alignment": {"status": "not_checked", "details": ""},
            "reply_to_risk": {"status": "not_checked", "details": ""}
        },
        "critical_issues": []
    }

    # Scoring weights (customizable)
    weights = {
        "dkim_fail": 35,
        "spf_fail": 35,
        "domain_mismatch": 25,
        "reply_to_risk": 15
    }

    # --- Validation Checks ---
    # DKIM Check
    if dkim_signature:
        if "v=1;" in dkim_signature.lower():
            results["detailed_checks"]["dkim"] = {
                "status": "valid", 
                "details": "Valid DKIM signature"
            }
        else:
            results["detailed_checks"]["dkim"] = {
                "status": "invalid",
                "details": "Invalid/malformed DKIM"
            }
            results["risk_score"] += weights["dkim_fail"]
            results["critical_issues"].append("DKIM_FAILURE")

    # SPF Check
    if received_spf:
        if "pass" in received_spf.lower():
            results["detailed_checks"]["spf"] = {
                "status": "valid",
                "details": "SPF authentication passed"
            }
        else:
            results["detailed_checks"]["spf"] = {
                "status": "invalid", 
                "details": "SPF authentication failed"
            }
            results["risk_score"] += weights["spf_fail"]
            results["critical_issues"].append("SPF_FAILURE")

    # Domain Alignment
    if from_header and return_path:
        from_domain = extract_domain(from_header)
        return_domain = extract_domain(return_path)
        
        if from_domain and return_domain:
            if from_domain != return_domain:
                results["detailed_checks"]["domain_alignment"] = {
                    "status": "invalid",
                    "details": f"From: {from_domain} ≠ Return-Path: {return_domain}"
                }
                results["risk_score"] += weights["domain_mismatch"]
            else:
                results["detailed_checks"]["domain_alignment"] = {
                    "status": "valid",
                    "details": f"Domains aligned: {from_domain}"
                }

    # Reply-To Risk
    if reply_to:
        reply_domain = extract_domain(reply_to)
        from_domain = extract_domain(from_header) if from_header else ""
        
        if reply_domain and from_domain and (reply_domain != from_domain):
            results["detailed_checks"]["reply_to_risk"] = {
                "status": "invalid",
                "details": f"Reply-To domain: {reply_domain}"
            }
            results["risk_score"] += weights["reply_to_risk"]

    # --- Risk Calculation ---
    results["risk_score"] = min(results["risk_score"], 100)
    
    # Determine risk level
    if results["risk_score"] >= 70:
        results["risk_level"] = "High"
    elif results["risk_score"] >= 40:
        results["risk_level"] = "Medium"
    else:
        results["risk_level"] = "Low"

    # Special case: Verified legitimate if full authentication
    if (results["detailed_checks"]["dkim"]["status"] == "valid" and
        results["detailed_checks"]["spf"]["status"] == "valid" and
        results["risk_score"] < 10):
        results["risk_level"] = "Verified Legitimate"

    # Final validity determination
    results["valid"] = results["risk_level"] in ("Low", "Verified Legitimate")
    
    return results