import numpy as np
from datetime import datetime
from typing import Dict, Any, List

def parse_timestamps(ts: List[str]) -> np.ndarray:
    if not ts:
        return np.array([])
    return np.array([datetime.fromisoformat(t) for t in ts])

def compute_iats(timestamps: np.ndarray) -> np.ndarray:
    if len(timestamps) < 2:
        return np.array([0])
    seconds = timestamps.astype('datetime64[ns]').astype(np.int64) / 1e9
    return np.diff(seconds)

def extract_cic_features(flow: Dict[str, Any]) -> Dict[str, Any]:
    fwd_lengths = np.array(flow.get("FwdPacketLengths") or [], dtype=float).reshape(-1)
    bwd_lengths = np.array(flow.get("BwdPacketLengths") or [], dtype=float).reshape(-1)
    all_lengths = np.concatenate((fwd_lengths, bwd_lengths))

    fwd_ts = np.array(parse_timestamps(flow.get("FwdTimestamps") or [])).reshape(-1)
    bwd_ts = np.array(parse_timestamps(flow.get("BwdTimestamps") or [])).reshape(-1)
    all_ts = np.sort(np.concatenate((fwd_ts, bwd_ts)))

    start_time = datetime.fromisoformat(flow["StartTime"])
    end_time = datetime.fromisoformat(flow["EndTime"])
    flow_duration = max((end_time - start_time).total_seconds(), 1e-6)

    flow_iat = compute_iats(all_ts)
    fwd_iat = compute_iats(fwd_ts)
    bwd_iat = compute_iats(bwd_ts)

    fin_flag = flow.get("FwdFIN", 0) + flow.get("BwdFIN", 0)
    syn_flag = flow.get("FwdSYN", 0) + flow.get("BwdSYN", 0)
    rst_flag = flow.get("FwdRST", 0) + flow.get("BwdRST", 0)
    psh_flag = flow.get("FwdPSH", 0) + flow.get("BwdPSH", 0)
    ack_flag = flow.get("FwdACK", 0) + flow.get("BwdACK", 0)
    urg_flag = flow.get("FwdURG", 0) + flow.get("BwdURG", 0)

    return {
        'Flow ID': flow["FlowID"],
        'Source IP': flow["SourceIP"],
        'Source Port': flow["SourcePort"],
        'Destination IP': flow["DestinationIP"],
        'Destination Port': flow["DestinationPort"],
        'Protocol': flow["Protocol"],
        'Flow Duration': flow_duration,
        'Total Fwd Packets': flow.get("TotalFwdPackets", 0),
        'Total Backward Packets': flow.get("TotalBwdPackets", 0),
        'Total Length of Fwd Packets': fwd_lengths.sum(),
        'Total Length of Bwd Packets': bwd_lengths.sum(),
        'Fwd Packet Length Max': fwd_lengths.max(initial=0),
        'Fwd Packet Length Min': fwd_lengths.min(initial=0),
        'Fwd Packet Length Mean': fwd_lengths.mean() if fwd_lengths.size > 0 else 0,
        'Fwd Packet Length Std': fwd_lengths.std() if fwd_lengths.size > 0 else 0,
        'Bwd Packet Length Max': bwd_lengths.max(initial=0),
        'Bwd Packet Length Min': bwd_lengths.min(initial=0),
        'Bwd Packet Length Mean': bwd_lengths.mean() if bwd_lengths.size > 0 else 0,
        'Bwd Packet Length Std': bwd_lengths.std() if bwd_lengths.size > 0 else 0,
        'Flow Bytes/s': (fwd_lengths.sum() + bwd_lengths.sum()) / flow_duration,
        'Flow Packets/s': (len(fwd_lengths) + len(bwd_lengths)) / flow_duration,
        'Flow IAT Mean': flow_iat.mean() if flow_iat.size > 0 else 0,
        'Flow IAT Std': flow_iat.std() if flow_iat.size > 0 else 0,
        'Flow IAT Max': flow_iat.max(initial=0),
        'Flow IAT Min': flow_iat.min(initial=0),
        'Fwd IAT Total': fwd_iat.sum(),
        'Fwd IAT Mean': fwd_iat.mean() if fwd_iat.size > 0 else 0,
        'Fwd IAT Std': fwd_iat.std() if fwd_iat.size > 0 else 0,
        'Fwd IAT Max': fwd_iat.max(initial=0),
        'Fwd IAT Min': fwd_iat.min(initial=0),
        'Bwd IAT Total': bwd_iat.sum(),
        'Bwd IAT Mean': bwd_iat.mean() if bwd_iat.size > 0 else 0,
        'Bwd IAT Std': bwd_iat.std() if bwd_iat.size > 0 else 0,
        'Bwd IAT Max': bwd_iat.max(initial=0),
        'Bwd IAT Min': bwd_iat.min(initial=0),
        'Fwd PSH Flags': flow.get("FwdPSH", 0),
        'Bwd PSH Flags': flow.get("BwdPSH", 0),
        'Fwd URG Flags': flow.get("FwdURG", 0),
        'Bwd URG Flags': flow.get("BwdURG", 0),
        'Fwd Header Length': flow.get("FwdHeaderLength", 0),
        'Bwd Header Length': flow.get("BwdHeaderLength", 0),
        'Fwd Packets/s': len(fwd_lengths) / flow_duration,
        'Bwd Packets/s': len(bwd_lengths) / flow_duration,
        'Min Packet Length': all_lengths.min(initial=0),
        'Max Packet Length': all_lengths.max(initial=0),
        'Packet Length Mean': all_lengths.mean() if all_lengths.size > 0 else 0,
        'Packet Length Std': all_lengths.std() if all_lengths.size > 0 else 0,
        'Packet Length Variance': all_lengths.var() if all_lengths.size > 0 else 0,
        'FIN Flag Count': fin_flag,
        'SYN Flag Count': syn_flag,
        'RST Flag Count': rst_flag,
        'PSH Flag Count': psh_flag,
        'ACK Flag Count': ack_flag,
        'URG Flag Count': urg_flag,
        'CWE Flag Count': 0,
        'ECE Flag Count': 0,
        'Down/Up Ratio': (len(bwd_lengths) / len(fwd_lengths)) if len(fwd_lengths) > 0 else 0,
        'Average Packet Size': all_lengths.mean() if all_lengths.size > 0 else 0,
        'Avg Fwd Segment Size': fwd_lengths.mean() if fwd_lengths.size > 0 else 0,
        'Avg Bwd Segment Size': bwd_lengths.mean() if bwd_lengths.size > 0 else 0,
        'Fwd Header Length.1': flow.get("FwdHeaderLength", 0),
        'Fwd Avg Bytes/Bulk': flow.get('FwdAvgBytesBulk', 0),
        'Fwd Avg Packets/Bulk': flow.get('FwdAvgPacketsBulk', 0),
        'Fwd Avg Bulk Rate': flow.get('FwdAvgBulkRate', 0),
        'Bwd Avg Bytes/Bulk': flow.get('BwdAvgBytesBulk', 0),
        'Bwd Avg Packets/Bulk': flow.get('BwdAvgPacketsBulk', 0),
        'Bwd Avg Bulk Rate': flow.get('BwdAvgBulkRate', 0),
        'Subflow Fwd Packets': len(fwd_lengths),
        'Subflow Fwd Bytes': fwd_lengths.sum(),
        'Subflow Bwd Packets': len(bwd_lengths),
        'Subflow Bwd Bytes': bwd_lengths.sum(),
        'act_data_pkt_fwd': len(fwd_lengths),
        'min_seg_size_forward': fwd_lengths.min(initial=0) if fwd_lengths.size > 0 else 0,
        'Active Mean': flow.get('ActiveMean', 0),
        'Active Std': flow.get('ActiveStd', 0),
        'Active Max': flow.get('ActiveMax', 0),
        'Active Min': flow.get('ActiveMin', 0),
        'Idle Mean': flow.get('IdleMean', 0),
        'Idle Std': flow.get('IdleStd', 0),
        'Idle Max': flow.get('IdleMax', 0),
        'Idle Min': flow.get('IdleMin', 0),
    }
