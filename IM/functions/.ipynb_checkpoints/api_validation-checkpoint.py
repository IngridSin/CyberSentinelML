from typing import Dict, Any
import re

def assess_email_risk(
    dkim_signature: str,
    received_spf: str,
    from_header: str = "",
    return_path: str = "",
    reply_to: str = "",
    message_id: str = ""
) -> Dict[str, Any]:
    """
    Assess email risk with nuanced scoring instead of binary validation.
    Returns risk score and detailed indicators.
    """
    
    def extract_domain(header: str) -> str:
        match = re.search(r'@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})', header or "")
        return match.group(1).lower() if match else ""

    # Initialize results with risk score (0-100)
    results = {
        "message_id": message_id,
        "risk_score": 0,
        "risk_level": "Low",
        "indicators": {
            "critical": [],
            "warnings": [],
            "info": []
        },
        "authentication_summary": {
            "dkim": False,
            "spf": False,
            "dmarc_aligned": False
        }
    }

    # Domain extraction
    from_domain = extract_domain(from_header)
    return_domain = extract_domain(return_path)
    reply_domain = extract_domain(reply_to)

    # --- Authentication Checks ---
    # DKIM Check
    dkim_valid = "v=1;" in (dkim_signature or "").lower()
    results["authentication_summary"]["dkim"] = dkim_valid
    if not dkim_valid:
        results["risk_score"] += 25
        results["indicators"]["critical"].append("Missing/Invalid DKIM")

    # SPF Check
    spf_valid = "pass" in (received_spf or "").lower()
    results["authentication_summary"]["spf"] = spf_valid
    if not spf_valid:
        results["risk_score"] += 25
        results["indicators"]["critical"].append("SPF validation failed")

    # DMARC Alignment (Simplified)
    results["authentication_summary"]["dmarc_aligned"] = dkim_valid and spf_valid
    if not results["authentication_summary"]["dmarc_aligned"]:
        results["risk_score"] += 15

    # --- Reputation Checks ---
    # Domain Alignment
    if from_domain and return_domain:
        if from_domain != return_domain:
            results["risk_score"] += 20
            results["indicators"]["warnings"].append(
                f"Domain mismatch: From({from_domain}) vs Return-Path({return_domain})"
            )
    else:
        results["risk_score"] += 10
        results["indicators"]["warnings"].append("Missing From/Return-Path header")

    # Reply-To Analysis
    if reply_domain and (reply_domain != from_domain):
        results["risk_score"] += 15
        results["indicators"]["warnings"].append(
            f"Suspicious Reply-To: {reply_domain}"
        )

    # --- Risk Level Calculation ---
    results["risk_score"] = min(results["risk_score"], 100)
    
    if results["risk_score"] >= 70:
        results["risk_level"] = "High"
    elif results["risk_score"] >= 40:
        results["risk_level"] = "Medium"
    else:
        results["risk_level"] = "Low"

    # --- Legitimate Email Safeguards ---
    if results["risk_score"] < 40 and not any(results["indicators"]["critical"]):
        results["risk_level"] = "Likely Legitimate"
        results["risk_score"] = max(results["risk_score"], 10)  # Prevent 0 scores

    return results
# def validate_email_headers(
#     dkim_signature: str,
#     received_spf: str,
#     from_header: str = "",
#     return_path: str = "",
#     reply_to: str = "",
#     message_id: str = ""
# ) -> Dict[str, Any]:
#     """
#     Validates email headers for phishing detection.
#     Handles missing From/Return-Path headers.
#     """
    
#     def extract_domain(header: str) -> str:
#         """Extracts domain from an email header."""
#         if not header:
#             return ""
#         # Match patterns like "user@domain.com" or "Name <user@domain.com>"
#         match = re.search(r'@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})', header)
#         return match.group(1).lower() if match else ""

#     # Initialize results
#     results = {
#         "message_id": message_id,
#         "checks": {
#             "dkim": {"valid": False, "details": "No DKIM signature"},
#             "spf": {"valid": False, "details": "SPF check failed"},
#             "dmarc_alignment": {"valid": False, "details": "DMARC alignment failed"},
#             "domain_alignment": {"valid": False, "details": "From/Return-Path mismatch"},
#             "reply_to_risk": {"valid": True, "details": "No risky Reply-To"},
#         },
#         "is_valid": False,
#         "warnings": []
#     }

#     # --- DKIM Check ---
#     if dkim_signature:
#         results["checks"]["dkim"]["valid"] = "v=1;" in dkim_signature.lower()
#         if results["checks"]["dkim"]["valid"]:
#             results["checks"]["dkim"]["details"] = "Valid DKIM signature"

#     # --- SPF Check ---
#     if received_spf:
#         results["checks"]["spf"]["valid"] = "pass" in received_spf.lower()
#         if results["checks"]["spf"]["valid"]:
#             results["checks"]["spf"]["details"] = "SPF passed"

#     # --- DMARC Alignment (Simplified) ---
#     # Requires SPF + DKIM alignment with From domain
#     if results["checks"]["dkim"]["valid"] and results["checks"]["spf"]["valid"]:
#         results["checks"]["dmarc_alignment"]["valid"] = True
#         results["checks"]["dmarc_alignment"]["details"] = "DMARC aligned"

#     # --- Domain Alignment ---
#     from_domain = extract_domain(from_header)
#     return_domain = extract_domain(return_path)
    
#     if from_domain and return_domain:
#         results["checks"]["domain_alignment"]["valid"] = (from_domain == return_domain)
#         if results["checks"]["domain_alignment"]["valid"]:
#             results["checks"]["domain_alignment"]["details"] = f"Domains aligned: {from_domain}"
#         else:
#             results["checks"]["domain_alignment"]["details"] = f"From: {from_domain} vs Return-Path: {return_domain}"
#     else:
#         results["checks"]["domain_alignment"]["details"] = "Missing From/Return-Path"
#         results["warnings"].append("Missing From or Return-Path header")

#     # --- Reply-To Check ---
#     reply_to_domain = extract_domain(reply_to)
#     if reply_to and reply_to_domain != from_domain:
#         results["checks"]["reply_to_risk"]["valid"] = False
#         details = f"Reply-To ({reply_to_domain}) ≠ From ({from_domain})" if from_domain else "Suspicious Reply-To"
#         results["checks"]["reply_to_risk"]["details"] = details

#     # --- Final Validity ---
#     results["is_valid"] = all([
#         results["checks"]["dkim"]["valid"],
#         results["checks"]["spf"]["valid"],
#         results["checks"]["dmarc_alignment"]["valid"],
#         results["checks"]["domain_alignment"]["valid"],
#         results["checks"]["reply_to_risk"]["valid"]
#     ])
    
#     return results


# def validate_email_headers(
#     dkim_signature: str,
#     received_spf: str,
#     from_header: str,
#     return_path: str,
#     message_id: str = ""
# ) -> Dict[str, Any]:
#     """
#     Validates email headers for phishing detection.
    
#     Args:
#         dkim_signature (str): DKIM-Signature header value
#         received_spf (str): Received-SPF header value
#         from_header (str): From header value
#         return_path (str): Return-Path header value
#         message_id (str): Message-ID header value (optional)
        
#     Returns:
#         Dict: Validation results with detailed checks
#     """
    
#     # Domain extraction helper
#     def extract_domain(header: str) -> str:
#         match = re.search(r'@([a-zA-Z0-9.-]+)', header)
#         return match.group(1).lower() if match else ""

#     # Perform checks
#     results = {
#         "message_id": message_id,
#         "checks": {
#             "dkim": {
#                 "valid": False,
#                 "details": "No valid DKIM signature found"
#             },
#             "spf": {
#                 "valid": False,
#                 "details": "SPF check failed"
#             },
#             "domain_alignment": {
#                 "valid": False,
#                 "details": "From/Return-Path domain mismatch"
#             }
#         },
#         "is_valid": False
#     }

#     # DKIM Validation
#     if dkim_signature:
#         results["checks"]["dkim"]["valid"] = "v=1;" in dkim_signature.lower()
#         if results["checks"]["dkim"]["valid"]:
#             results["checks"]["dkim"]["details"] = "Valid DKIM signature found"

#     # SPF Validation
#     if received_spf:
#         results["checks"]["spf"]["valid"] = "pass" in received_spf.lower()
#         if results["checks"]["spf"]["valid"]:
#             results["checks"]["spf"]["details"] = "SPF check passed"

#     # Domain Alignment Check
#     from_domain = extract_domain(from_header)
#     return_domain = extract_domain(return_path)
    
#     if from_domain and return_domain:
#         results["checks"]["domain_alignment"]["valid"] = (from_domain == return_domain)
#         if results["checks"]["domain_alignment"]["valid"]:
#             results["checks"]["domain_alignment"]["details"] = \
#                 f"Domains aligned: {from_domain}"
#         else:
#             results["checks"]["domain_alignment"]["details"] = \
#                 f"From: {from_domain} vs Return-Path: {return_domain}"

#     # Final validation
#     results["is_valid"] = all([
#         results["checks"]["dkim"]["valid"],
#         results["checks"]["spf"]["valid"],
#         results["checks"]["domain_alignment"]["valid"]
#     ])
    
#     return results