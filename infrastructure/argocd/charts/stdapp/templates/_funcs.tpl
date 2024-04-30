{{- define "deployment.name"}}
{{- required "name must be set" .Values.name }}
{{- end }}

{{- define "deployment.labels.matchLabels" }}
app: {{ include "deployment.name" . }}
app.kubernetes.io/name: {{ include "deployment.name" . }}
{{- end }}

{{- define "deployment.labels.standard" }}
{{- include "deployment.labels.matchLabels" . }}
helm.sh/chart: {{ .Chart.Name }}
app.kubernetes.io/managed-by: helm
app.kubernetes.io/version: {{ required "version must be set" .Values.version }}
version: {{ required "version must be set" .Values.version }}
{{- end }}
