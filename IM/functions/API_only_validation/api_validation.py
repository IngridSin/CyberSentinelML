from typing import Dict, Any
import re

def validate_email_headers(
    dkim_signature: str,
    received_spf: str,
    from_header: str = "",
    return_path: str = "",
    reply_to: str = "",
    message_id: str = ""
) -> Dict[str, Any]:
    """
    Validates email headers for phishing detection.
    Handles missing From/Return-Path headers.
    """
    
    def extract_domain(header: str) -> str:
        """Extracts domain from an email header."""
        if not header:
            return ""
        match = re.search(r'@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})', header)
        return match.group(1).lower() if match else ""

    # Initialize results
    results = {
        "message_id": message_id,
        "checks": {
            "dkim": {"valid": False, "details": "No DKIM signature"},
            "spf": {"valid": False, "details": "SPF check failed"},
            "dmarc_alignment": {"valid": False, "details": "DMARC alignment failed"},
            "domain_alignment": {"valid": False, "details": "From/Return-Path mismatch"},
            "reply_to_risk": {"valid": True, "details": "No risky Reply-To"},
        },
        "is_valid": False,
        "warnings": []
    }

    # --- DKIM Check ---
    if dkim_signature:
        results["checks"]["dkim"]["valid"] = "v=1;" in dkim_signature.lower()
        if results["checks"]["dkim"]["valid"]:
            results["checks"]["dkim"]["details"] = "Valid DKIM signature"

    # --- SPF Check ---
    if received_spf:
        results["checks"]["spf"]["valid"] = "pass" in received_spf.lower()
        if results["checks"]["spf"]["valid"]:
            results["checks"]["spf"]["details"] = "SPF passed"

    # --- DMARC Alignment (Simplified) ---
    # Requires SPF + DKIM alignment with From domain
    if results["checks"]["dkim"]["valid"] and results["checks"]["spf"]["valid"]:
        results["checks"]["dmarc_alignment"]["valid"] = True
        results["checks"]["dmarc_alignment"]["details"] = "DMARC aligned"

    # --- Domain Alignment ---
    from_domain = extract_domain(from_header)
    return_domain = extract_domain(return_path)
    
    if from_domain and return_domain:
        results["checks"]["domain_alignment"]["valid"] = (from_domain == return_domain)
        if results["checks"]["domain_alignment"]["valid"]:
            results["checks"]["domain_alignment"]["details"] = f"Domains aligned: {from_domain}"
        else:
            results["checks"]["domain_alignment"]["details"] = f"From: {from_domain} vs Return-Path: {return_domain}"
    else:
        results["checks"]["domain_alignment"]["details"] = "Missing From/Return-Path"
        results["warnings"].append("Missing From or Return-Path header")

    # --- Reply-To Check ---
    reply_to_domain = extract_domain(reply_to)
    if reply_to and reply_to_domain != from_domain:
        results["checks"]["reply_to_risk"]["valid"] = False
        details = f"Reply-To ({reply_to_domain}) ≠ From ({from_domain})" if from_domain else "Suspicious Reply-To"
        results["checks"]["reply_to_risk"]["details"] = details

    # --- Final Validity ---
    results["is_valid"] = all([
        results["checks"]["dkim"]["valid"],
        results["checks"]["spf"]["valid"],
        results["checks"]["dmarc_alignment"]["valid"],
        results["checks"]["domain_alignment"]["valid"],
        results["checks"]["reply_to_risk"]["valid"]
    ])
    
    return results

