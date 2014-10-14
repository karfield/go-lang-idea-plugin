package ro.redeul.google.go.lang.psi.resolve;

import com.intellij.psi.impl.source.resolve.ResolveCache.AbstractResolver;
import org.jetbrains.annotations.NotNull;
import ro.redeul.google.go.lang.psi.resolve.references.Reference;

public class ResolveCacheResolvers {

    public static <
            Ref extends Reference<?, ?, Solver, Ref>,
            Solver extends RefSolver<Ref, Solver>
            > AbstractResolver<Ref, GoResolveResult> makeDefault() {
        return new AbstractResolver<Ref, GoResolveResult>() {
            @Override
            public GoResolveResult resolve(@NotNull Ref reference, boolean incompleteCode) {
                Solver processor = reference.newSolver();

                reference.walkSolver(processor);

                return GoResolveResult.fromElement(processor.getChildDeclaration());
            }
        };
    }
}
