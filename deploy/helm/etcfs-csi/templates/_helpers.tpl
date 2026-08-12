{{- define "etcfs-csi.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "etcfs-csi.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "etcfs-csi.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "etcfs-csi.labels" -}}
app.kubernetes.io/name: {{ include "etcfs-csi.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "etcfs-csi.image" -}}
{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}
{{- end -}}

{{/*
etcd client flags, shared by every container that talks to etcd.
*/}}
{{- define "etcfs-csi.etcdArgs" -}}
- --etcd-endpoints={{ required "etcd.endpoints is required: point the driver at the EtcFS cluster's etcd, not the Kubernetes control plane's" .Values.etcd.endpoints }}
{{- if .Values.etcd.tlsSecretName }}
- --etcd-ca=/etc/etcfs/etcd-tls/ca.crt
- --etcd-cert=/etc/etcfs/etcd-tls/tls.crt
- --etcd-key=/etc/etcfs/etcd-tls/tls.key
{{- end }}
{{- end -}}
