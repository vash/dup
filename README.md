# Dup - Kubectl Plugin for effective duplication of Kubernetes Resources

This plugin is designed for on-the-fly duplication of Kubernetes resources.
It focuses on providing a convenient way to edit resources before duplication,
with a specific emphasis on Pods to create a fine-tuned resource quickly.
This tool can be used for debugging running containers without them crashing,
and simplifying the administration and general interaction with Kubernetes clusters.

## Installation

### Using Krew (Kubectl Plugin Manager) # SOON

```bash
kubectl krew install dup
```

### Manual Installation

Download the latest release from the GitHub releases page.
Extract the binary from the archive.
Move the binary to a directory in your system's PATH.

### Usage

```bash
kubectl dup [options] <resource-type,resource-type-2> <resource-name> <generated-resource-name-prefix>
```

## Examples

```bash
# duplicate a pod of deployment "my-deployment" without opening edit window
kubectl dup deployment my-deployment -pk

# duplicate a specific pod, regardless of what it belongs to,
# and disable liveness probes on it.
kubectl dup pod my-pod -d
```

## Features

- Resource Duplication: Duplicate Kubernetes resources with ease.
- Interactive Editing: Open resources for manual editing before applying changes.
- Pod Focused: Designed to streamline debugging and administration of Pods.

## Options

-d, --disable-probes Automatically disable readiness and liveness probes for duplicated pods.
-h, --help: Display help information.
-p, --pod: Duplicate pod of complex objects, supported objects: 'StatefulSet','Deployment','CronJob','Job'.
-k, --skip-edit: Skip editing duplicated resource before creation

## Contributing

We welcome contributions! If you have any ideas, enhancements, or bug fixes,
feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
